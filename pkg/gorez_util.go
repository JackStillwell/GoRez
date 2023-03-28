package gorez

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	i "github.com/JackStillwell/GoRez/pkg/interfaces"

	authI "github.com/JackStillwell/GoRez/internal/auth_service/interfaces"

	requestI "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	requestM "github.com/JackStillwell/GoRez/internal/request_service/models"
	requestU "github.com/JackStillwell/GoRez/internal/request_service/utilities"

	sessionI "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
	sessionM "github.com/JackStillwell/GoRez/internal/session_service/models"
)

type gorezUtil struct {
	authSvc authI.AuthService
	rqstSvc requestI.RequestService
	sesnSvc sessionI.SessionService
}

func NewGorezUtil(aS authI.AuthService, rS requestI.RequestService,
	sS sessionI.SessionService) i.GorezUtil {
	return &gorezUtil{
		authSvc: aS,
		rqstSvc: rS,
		sesnSvc: sS,
	}
}

func (g *gorezUtil) BulkAsyncSessionRequest(requestBuilders []func(*sessionM.Session) *requestM.Request,
) ([][]byte, []error) {
	numRequests := len(requestBuilders)
	responseChan := make(chan *requestM.RequestResponse, numRequests)
	uIDSessionMap := make(map[*uuid.UUID]*sessionM.Session, numRequests)

	// NOTE: this is async so the reservation and release of sessions is possible, but the func
	// return depends upon responses being completed.
	go func() {

		for _, rB := range requestBuilders {
			sessChan := make(chan *sessionM.Session, 1)
			g.sesnSvc.ReserveSession(1, sessChan)

			s := <-sessChan // NOTE: will wait here until session recieved

			// Request constructed with session here
			r := rB(s)

			uID := uuid.New()
			r.Id = &uID
			uIDSessionMap[&uID] = s

			g.rqstSvc.MakeRequest(r)
			g.rqstSvc.GetResponse(&uID, responseChan)
		}
	}()

	responses := make([][]byte, numRequests)
	errs := make([]error, numRequests)
	for i := 0; i < numRequests; i++ {
		resp := <-responseChan
		log.Println("response received:", resp)
		if resp.Err != nil {
			if strings.Contains(resp.Err.Error(), "session") {
				g.sesnSvc.BadSession([]*sessionM.Session{uIDSessionMap[resp.Id]})
			} else {
				g.sesnSvc.ReleaseSession([]*sessionM.Session{uIDSessionMap[resp.Id]})
			}
			errs[i] = errors.Wrap(resp.Err, "request")
			continue
		}

		g.sesnSvc.ReleaseSession([]*sessionM.Session{uIDSessionMap[resp.Id]})

		responses[i] = resp.Resp

		log.Println("responses at", i, "equals", resp.Resp)
	}

	log.Println("responses", responses)

	return responses, errs
}

func (g *gorezUtil) MultiRequest(requestArgs []string, baseURL, method string,
) ([][]byte, []error) {
	requestBuilders := make([]func(*sessionM.Session) *requestM.Request, len(requestArgs))

	for i, arg := range requestArgs {
		requestBuilders[i] = func(s *sessionM.Session) *requestM.Request {
			return &requestM.Request{
				JITArgs: []interface{}{
					baseURL,
					g.authSvc.GetID(),
					method,
					s.Key,
					g.authSvc.GetTimestamp,
					g.authSvc.GetSignature,
					arg,
				},
				JITBuild: requestU.JITBase,
			}
		}
	}

	rawObjs, errs := g.BulkAsyncSessionRequest(requestBuilders)

	return rawObjs, errs
}

func (g *gorezUtil) SingleRequest(r requestM.Request, unmarshalTo interface{}) error {
	log.Println("reserving session for single request")
	sesnChan := make(chan *sessionM.Session, 1)
	g.sesnSvc.ReserveSession(1, sesnChan)
	s := <-sesnChan
	log.Println("session reserved for single request")

	sessions := []*sessionM.Session{s}
	defer g.sesnSvc.ReleaseSession(sessions)

	r.JITArgs[1] = g.authSvc.GetID()
	r.JITArgs[3] = s.Key
	r.JITArgs[4] = g.authSvc.GetTimestamp
	r.JITArgs[5] = g.authSvc.GetSignature

	log.Println("making single request")
	resp := g.rqstSvc.Request(&r)
	log.Println("single response received")

	if resp.Err != nil {
		return errors.Wrap(resp.Err, "requesting response")
	}

	err := json.Unmarshal(resp.Resp, &unmarshalTo)
	if err != nil {
		return errors.Wrap(err, "unmarshaling response")
	}

	log.Println("single response unmarshaled")

	return nil
}
