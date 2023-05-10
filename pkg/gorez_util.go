package gorez

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

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

	sessionMapLock := sync.RWMutex{}
	uIDSessionMap := make(map[*uuid.UUID]*sessionM.Session, numRequests)
	getSession := func(uID *uuid.UUID) *sessionM.Session {
		log.Println("waiting for session map lock")
		sessionMapLock.Lock()
		log.Println("session map lock acquired")
		retVal := uIDSessionMap[uID]
		log.Println("releasing session map lock")
		sessionMapLock.Unlock()
		return retVal
	}
	setSession := func(uID *uuid.UUID, sess *sessionM.Session) {
		log.Println("waiting for session map lock")
		sessionMapLock.RLock()
		log.Println("session map lock acquired")
		uIDSessionMap[uID] = sess
		log.Println("releasing session map lock")
		sessionMapLock.RUnlock()
	}

	responseMapLock := sync.RWMutex{}
	uIDResponseIdxMap := make(map[*uuid.UUID]int, numRequests)
	getResponseIdx := func(uID *uuid.UUID) int {
		log.Println("waiting for responseIdx map lock")
		responseMapLock.RLock()
		log.Println("responseIdx map lock acquired")
		retVal := uIDResponseIdxMap[uID]
		log.Println("releasing responseIdx map lock")
		responseMapLock.RUnlock()
		return retVal
	}
	setResponseIdx := func(uID *uuid.UUID, idx int) {
		log.Println("waiting for responseIdx map lock")
		responseMapLock.Lock()
		log.Println("responseIdx map lock acquired")
		uIDResponseIdxMap[uID] = idx
		log.Println("releasing responseIdx map lock")
		responseMapLock.Unlock()
	}

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

			setSession(&uID, s)
			setResponseIdx(&uID, i)

			g.rqstSvc.MakeRequest(r)
			uIDs <- &uID
		}
	}()

	responses := make([][]byte, numRequests)
	errs := make([]error, numRequests)
	for i := 0; i < numRequests; i++ {
		uID := <-uIDs

		idx := getResponseIdx(uID)

		resp := g.rqstSvc.GetResponse(uID)
		sess := getSession(resp.Id)
		if resp.Err != nil {
			if strings.Contains(resp.Err.Error(), "session") {
				g.sesnSvc.BadSession([]*sessionM.Session{sess})
			} else {
				g.sesnSvc.ReleaseSession([]*sessionM.Session{sess})
			}
			errs[idx] = fmt.Errorf("request: %w", resp.Err)
			continue
		}

		g.sesnSvc.ReleaseSession([]*sessionM.Session{sess})

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
