package gorez

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
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
		// WARN: this is required for testing recovery. Uncomment if debugging test failures.
		defer ginkgo.GinkgoRecover()

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
	}

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
	sesnChan := make(chan *sessionM.Session, 1)
	g.sesnSvc.ReserveSession(1, sesnChan)
	s := <-sesnChan

	sessions := []*sessionM.Session{s}
	defer g.sesnSvc.ReleaseSession(sessions)

	r.JITArgs[1] = g.authSvc.GetID()
	r.JITArgs[3] = s.Key
	r.JITArgs[4] = g.authSvc.GetTimestamp
	r.JITArgs[5] = g.authSvc.GetSignature

	resp := g.rqstSvc.Request(&r)

	if resp.Err != nil {
		return errors.Wrap(resp.Err, "requesting response")
	}

	err := json.Unmarshal(resp.Resp, &unmarshalTo)
	if err != nil {
		return errors.Wrap(err, "unmarshaling response")
	}

	return nil
}
