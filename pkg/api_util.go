package gorez

import (
	"encoding/json"

	authService "github.com/JackStillwell/GoRez/internal/auth_service/interfaces"
	requestService "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	rSM "github.com/JackStillwell/GoRez/internal/request_service/models"
	requestUtils "github.com/JackStillwell/GoRez/internal/request_service/utilities"
	sessionService "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
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

func (a *apiUtil) CreateSession() (*m.Session, error) {
	uID := uuid.New()
	r := rSM.Request{
		Id: &uID,
		JITArgs: []interface{}{
			"http://api.smitegame.com/smiteapi.svc/createsessionjson",
			a.authSvc.GetID(),
			"createsession",
			"",
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

	session := &m.Session{}
	err := json.Unmarshal(resp.Resp, session)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshall response")
	}

	return session, nil
}

// TODO: make bulk by taking a list and returning a list
func (a *apiUtil) TestSession(s *m.Session) (string, error) {
	uID := uuid.New()
	r := rSM.Request{
		Id: &uID,
		JITArgs: []interface{}{
			hRConst.SmiteURLBase + hRConst.TestSession + "json",
			a.authSvc.GetID(),
			hRConst.TestSession,
			s,
			a.authSvc.GetTimestamp,
			a.authSvc.GetSignature,
			"",
		},
		JITBuild: requestUtils.JITBase,
	}

	a.rqstSvc.MakeRequest(&r)
	respChan := make(chan *rSM.RequestResponse, 1)
	a.rqstSvc.GetResponse(&uID, respChan)
	resp := <-respChan
	if resp.Err != nil {
		return "", errors.Wrap(resp.Err, "request")
	}

	return string(resp.Resp), nil
}

func (a *apiUtil) GetDataUsed() (*m.UsageInfo, error) {
	sessions, err := a.sesnSvc.ReserveSession(1)
	if err != nil {
		return nil, errors.Wrap(err, "reserving session")
	}

	defer a.sesnSvc.ReleaseSession(sessions)
	s := sessions[0]

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
	err = json.Unmarshal(resp.Resp, uI)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling response")
	}

	return uI, nil
}
