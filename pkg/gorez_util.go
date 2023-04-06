package gorez

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"

	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"

	authI "github.com/JackStillwell/GoRez/internal/auth/interfaces"

	requestI "github.com/JackStillwell/GoRez/internal/request/interfaces"
	requestM "github.com/JackStillwell/GoRez/internal/request/models"

	sessionI "github.com/JackStillwell/GoRez/internal/session/interfaces"
	sessionM "github.com/JackStillwell/GoRez/internal/session/models"
)

type gorezUtil struct {
	authSvc authI.Service
	rqstSvc requestI.Service
	sesnSvc sessionI.Service
}

func NewGorezUtil(aS authI.Service, rS requestI.Service,
	sS sessionI.Service) i.GorezUtil {
	return &gorezUtil{
		authSvc: aS,
		rqstSvc: rS,
		sesnSvc: sS,
	}
}

func (g *gorezUtil) BulkAsyncSessionRequest(requestBuilders []func(*sessionM.Session) *requestM.Request,
) ([][]byte, []error) {
	numRequests := len(requestBuilders)
	uIDs := make(chan *uuid.UUID, numRequests)
	uIDSessionMap := make(map[*uuid.UUID]*sessionM.Session, numRequests)
	uIDResponseIdxMap := make(map[*uuid.UUID]int, numRequests)

	// NOTE: this is async so the reservation and release of sessions is possible, but the func
	// return depends upon responses being completed.
	go func() {
		for i, rB := range requestBuilders {
			sessChan := make(chan *sessionM.Session, 1)
			g.sesnSvc.ReserveSession(1, sessChan)

			s := <-sessChan // NOTE: will wait here until session recieved

			// Request constructed with session here
			r := rB(s)

			uID := uuid.New()
			r.Id = &uID
			uIDSessionMap[&uID] = s
			uIDResponseIdxMap[&uID] = i

			g.rqstSvc.MakeRequest(r)
			uIDs <- &uID
		}
	}()

	responses := make([][]byte, numRequests)
	errs := make([]error, numRequests)
	for i := 0; i < numRequests; i++ {
		uID := <-uIDs
		idx := uIDResponseIdxMap[uID]
		resp := g.rqstSvc.GetResponse(uID)
		if resp.Err != nil {
			if strings.Contains(resp.Err.Error(), "session") {
				g.sesnSvc.BadSession([]*sessionM.Session{uIDSessionMap[resp.Id]})
			} else {
				g.sesnSvc.ReleaseSession([]*sessionM.Session{uIDSessionMap[resp.Id]})
			}
			errs[idx] = fmt.Errorf("request: %w", resp.Err)
			continue
		}

		g.sesnSvc.ReleaseSession([]*sessionM.Session{uIDSessionMap[resp.Id]})

		responses[idx] = resp.Resp
	}

	return responses, errs
}

func (g *gorezUtil) MultiRequest(requestArgs []string, baseURL, method string,
) ([][]byte, []error) {
	requestBuilders := make([]func(*sessionM.Session) *requestM.Request, len(requestArgs))

	for i, arg := range requestArgs {
		requestBuilders[i] = func(s *sessionM.Session) *requestM.Request {
			return &requestM.Request{
				JITFunc: HiRezJIT(
					baseURL,
					g.authSvc.GetID(),
					method,
					s.Key,
					g.authSvc.GetTimestamp,
					g.authSvc.GetSignature,
					arg,
				),
			}
		}
	}

	return g.BulkAsyncSessionRequest(requestBuilders)
}

func (g *gorezUtil) SingleRequest(url, endpoint, endpointArgs string) ([]byte, error) {
	log.Println("reserving session for single request")
	sesnChan := make(chan *sessionM.Session, 1)
	g.sesnSvc.ReserveSession(1, sesnChan)
	s := <-sesnChan
	log.Println("session reserved for single request")

	sessions := []*sessionM.Session{s}
	defer g.sesnSvc.ReleaseSession(sessions)

	r := requestM.Request{
		JITFunc: HiRezJIT(
			url, g.authSvc.GetID(), endpoint, s.Key, g.authSvc.GetTimestamp, g.authSvc.GetSignature,
			endpointArgs,
		),
	}

	log.Println("making single request")
	resp := g.rqstSvc.Request(&r)
	log.Println("single response received")

	if resp.Err != nil {
		return nil, fmt.Errorf("requesting response: %w", resp.Err)
	}

	if resp.Resp[0] == byte('[') {
		retMsgs := []m.RetMsg{}
		err := json.Unmarshal(resp.Resp, &retMsgs)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling response ret msg: %w", err)
		}

		log.Println("single response unmarshaled")

		for i, retMsg := range retMsgs {
			if retMsg.Msg != nil && *retMsg.Msg != "" {
				return nil, fmt.Errorf("ret_msg %d: %s", i, *retMsg.Msg)
			}
		}

	} else {
		retMsg := m.RetMsg{}
		err := json.Unmarshal(resp.Resp, &retMsg)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling response ret msg: %w", err)
		}

		log.Println("single response unmarshaled")

		if retMsg.Msg != nil && *retMsg.Msg != "" {
			return nil, fmt.Errorf("ret_msg: %s", *retMsg.Msg)
		}
	}

	return resp.Resp, nil
}
