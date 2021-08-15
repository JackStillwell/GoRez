package gorez

import (
	"encoding/json"

	authService "github.com/JackStillwell/GoRez/internal/auth_service/interfaces"
	requestService "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	rSM "github.com/JackStillwell/GoRez/internal/request_service/models"
	requestUtils "github.com/JackStillwell/GoRez/internal/request_service/utilities"
	sessionService "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
	sSM "github.com/JackStillwell/GoRez/internal/session_service/models"
	hRConst "github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type apiUtil struct {
	authSvc authService.AuthService
	rqstSvc requestService.RequestService
	sesnSvc sessionService.SessionService
}

func NewAPIUtil(
	rS requestService.RequestService,
	aS authService.AuthService,
	sS sessionService.SessionService,
) i.APIUtil {
	return &apiUtil{
		rqstSvc: rS,
		authSvc: aS,
		sesnSvc: sS,
	}
}

func (a *apiUtil) CreateSession(numSessions int) ([]*m.Session, []error) {
	r := rSM.Request{
		JITArgs: []interface{}{
			hRConst.SmiteURLBase + hRConst.CreateSession + "json",
			a.authSvc.GetID(),
			hRConst.CreateSession,
			"",
			a.authSvc.GetTimestamp,
			a.authSvc.GetSignature,
			"",
		},
		JITBuild: requestUtils.JITBase,
	}

	uIDs := make([]*uuid.UUID, numSessions)
	for i := 0; i < numSessions; i++ {
		uID := uuid.New()
		r.Id = &uID
		a.rqstSvc.MakeRequest(&r)
		uIDs = append(uIDs, &uID)
	}

	responseChan := make(chan *rSM.RequestResponse, numSessions)
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
			errs = append(errs, errors.Wrap(err, "unmarshall response"))
			continue
		}

		sessions = append(sessions, session)
	}

	return sessions, errs
}

func (a *apiUtil) TestSession(s []*m.Session) ([]*string, []error) {
	r := rSM.Request{
		JITArgs: []interface{}{
			hRConst.SmiteURLBase + hRConst.TestSession + "json",
			a.authSvc.GetID(),
			hRConst.TestSession,
			"",
			a.authSvc.GetTimestamp,
			a.authSvc.GetSignature,
			"",
		},
		JITBuild: requestUtils.JITBase,
	}

	uIDs := make([]*uuid.UUID, len(s))
	for i := 0; i < len(s); i++ {
		uID := uuid.New()
		r.Id = &uID
		a.rqstSvc.MakeRequest(&r)
		uIDs = append(uIDs, &uID)
	}

	responseChan := make(chan *rSM.RequestResponse, len(s))
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
	sesnChan := make(chan *sSM.Session, 1)
	a.sesnSvc.ReserveSession(1, sesnChan)
	s := <-sesnChan

	sessions := []*sSM.Session{s}
	defer a.sesnSvc.ReleaseSession(sessions)

	uID := uuid.New()
	r := rSM.Request{
		Id: &uID,
		JITArgs: []interface{}{
			hRConst.SmiteURLBase + hRConst.GetDataUsed + "json",
			a.authSvc.GetID(),
			hRConst.GetDataUsed,
			s.Key,
			a.authSvc.GetTimestamp,
			a.authSvc.GetSignature,
			"",
		},
		JITBuild: requestUtils.JITBase,
	}

	resp := a.rqstSvc.Request(&r)
	if resp.Err != nil {
		return nil, errors.Wrap(resp.Err, "request")
	}

	uI := &m.UsageInfo{}
	err := json.Unmarshal(resp.Resp, uI)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling response")
	}

	return uI, nil
}
