package gorez

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	c "github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"

	authI "github.com/JackStillwell/GoRez/internal/auth_service/interfaces"

	requestI "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	requestM "github.com/JackStillwell/GoRez/internal/request_service/models"
	requestU "github.com/JackStillwell/GoRez/internal/request_service/utilities"

	sessionI "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
	sessionM "github.com/JackStillwell/GoRez/internal/session_service/models"
)

type apiUtil struct {
	authSvc authI.AuthService
	rqstSvc requestI.RequestService
	sesnSvc sessionI.SessionService
}

func NewAPIUtil(
	rS requestI.RequestService,
	aS authI.AuthService,
	sS sessionI.SessionService,
) i.APIUtil {
	return &apiUtil{
		rqstSvc: rS,
		authSvc: aS,
		sesnSvc: sS,
	}
}

func (a *apiUtil) CreateSession(numSessions int) ([]*m.Session, []error) {
	r := requestM.Request{
		JITArgs: []interface{}{
			c.SmiteURLBase + c.CreateSession + "json",
			a.authSvc.GetID(),
			c.CreateSession,
			"",
			a.authSvc.GetTimestamp,
			a.authSvc.GetSignature,
			"",
		},
		JITBuild: requestU.JITBase,
	}

	uIDs := make([]*uuid.UUID, numSessions)
	for i := 0; i < numSessions; i++ {
		uID := uuid.New()
		r.Id = &uID
		a.rqstSvc.MakeRequest(&r)
		uIDs = append(uIDs, &uID)
	}

	responseChan := make(chan *requestM.RequestResponse, numSessions)
	for i := 0; i < numSessions; i++ {
		a.rqstSvc.GetResponse(uIDs[i], responseChan)
	}

	sessions := make([]*m.Session, 0, numSessions)
	errs := make([]error, 0, numSessions)
	for i := 0; i < numSessions; i++ {
		resp := <-responseChan
		if resp.Err != nil {
			errs = append(errs, errors.Wrap(resp.Err, "request"))
			continue
		}

		session := &m.Session{}
		err := json.Unmarshal(resp.Resp, session)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "unmarshal response"))
			continue
		}

		sessions = append(sessions, session)
	}

	return sessions, errs
}

func (a *apiUtil) TestSession(s []*m.Session) ([]*string, []error) {
	r := requestM.Request{
		JITArgs: []interface{}{
			c.SmiteURLBase + c.TestSession + "json",
			a.authSvc.GetID(),
			c.TestSession,
			"",
			a.authSvc.GetTimestamp,
			a.authSvc.GetSignature,
			"",
		},
		JITBuild: requestU.JITBase,
	}

	uIDs := make([]*uuid.UUID, len(s))
	for i := 0; i < len(s); i++ {
		uID := uuid.New()
		r.Id = &uID
		a.rqstSvc.MakeRequest(&r)
		uIDs = append(uIDs, &uID)
	}

	responseChan := make(chan *requestM.RequestResponse, len(s))
	for i := 0; i < len(s); i++ {
		a.rqstSvc.GetResponse(uIDs[i], responseChan)
	}

	responses := make([]*string, 0, len(s))
	errs := make([]error, 0, len(s))
	for i := 0; i < len(s); i++ {
		resp := <-responseChan
		if resp.Err != nil {
			errs = append(errs, errors.Wrap(resp.Err, "request"))
			continue
		}

		responseString := string(resp.Resp)
		responses = append(responses, &responseString)
	}

	return responses, errs
}

// NOTE: can only do one at a time, so no need for bulk concurrency
func (a *apiUtil) GetDataUsed() (*m.UsageInfo, error) {
	sesnChan := make(chan *sessionM.Session, 1)
	a.sesnSvc.ReserveSession(1, sesnChan)
	s := <-sesnChan

	sessions := []*sessionM.Session{s}
	defer a.sesnSvc.ReleaseSession(sessions)

	uID := uuid.New()
	r := requestM.Request{
		Id: &uID,
		JITArgs: []interface{}{
			c.SmiteURLBase + c.GetDataUsed + "json",
			a.authSvc.GetID(),
			c.GetDataUsed,
			s.Key,
			a.authSvc.GetTimestamp,
			a.authSvc.GetSignature,
			"",
		},
		JITBuild: requestU.JITBase,
	}

	resp := a.rqstSvc.Request(&r)
	if resp.Err != nil {
		return nil, errors.Wrap(resp.Err, "request")
	}

	uI := &m.UsageInfo{}
	err := json.Unmarshal(resp.Resp, uI)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshaling response")
	}

	return uI, nil
}
