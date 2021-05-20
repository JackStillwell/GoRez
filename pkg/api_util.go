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
	aS authService.AuthService
	rS requestService.RequestService
	sS sessionService.SessionService
}

func NewAPIUtil(
	rS requestService.RequestService,
	aS authService.AuthService,
	sS sessionService.SessionService,
) i.APIUtil {
	return &apiUtil{
		rS: rS,
		aS: aS,
		sS: sS,
	}
}

func (a *apiUtil) CreateSession() (*m.Session, error) {
	uID := uuid.New()
	r := rSM.Request{
		Id: &uID,
		JITArgs: []interface{}{
			"http://api.smitegame.com/smiteapi.svc/createsessionjson",
			a.aS.GetID(),
			"createsession",
			"",
			a.aS.GetTimestamp,
			a.aS.GetSignature,
			"",
		},
		JITBuild: requestUtils.JITBase,
	}

	resp := a.rS.Request(&r)
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

func (a *apiUtil) TestSession(s *m.Session) (string, error) {
	uID := uuid.New()
	r := rSM.Request{
		Id: &uID,
		JITArgs: []interface{}{
			hRConst.SmiteURLBase + hRConst.TestSession + "json",
			a.aS.GetID(),
			hRConst.TestSession,
			s,
			a.aS.GetTimestamp,
			a.aS.GetSignature,
			"",
		},
		JITBuild: requestUtils.JITBase,
	}

	a.rS.MakeRequest(&r)
	resp := a.rS.GetResponse(&uID)
	if resp.Err != nil {
		return "", errors.Wrap(resp.Err, "request")
	}

	return string(resp.Resp), nil
}

func (a *apiUtil) GetDataUsed() (*m.UsageInfo, error) {
	sessions, err := a.sS.ReserveSession(1)
	if err != nil {
		return nil, errors.Wrap(err, "reserving session")
	}

	defer a.sS.ReleaseSession(sessions)
	s := sessions[0]

	uID := uuid.New()
	r := rSM.Request{
		Id: &uID,
		JITArgs: []interface{}{
			hRConst.SmiteURLBase + hRConst.GetDataUsed + "json",
			a.aS.GetID(),
			hRConst.GetDataUsed,
			s.Key,
			a.aS.GetTimestamp,
			a.aS.GetSignature,
			"",
		},
		JITBuild: requestUtils.JITBase,
	}

	resp := a.rS.Request(&r)
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
