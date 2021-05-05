package gorez

import (
	"encoding/json"
	"fmt"
	"time"

	authService "github.com/JackStillwell/GoRez/internal/auth_service/interfaces"
	requestService "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	rSM "github.com/JackStillwell/GoRez/internal/request_service/models"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"
	"github.com/google/uuid"
)

type apiUtil struct {
	aS authService.AuthService
	rS requestService.RequestService
}

func NewAPIUtil(rS requestService.RequestService, aS authService.AuthService) i.APIUtil {
	return &apiUtil{
		rS: rS,
		aS: aS,
	}
}

func (a *apiUtil) CreateSession() *m.Session {
	jitFunc := func(args []interface{}) (string, error) {
		baseURL := "http://api.smitegame.com/smiteapi.svc/createsessionjson"
		t := time.Now().UTC()
		tS := a.aS.GetTimestamp(t)
		return fmt.Sprintf(
			"%s/%s/%s/%s",
			baseURL,
			a.aS.GetID(),
			a.aS.GetSignature("createsession", tS),
			tS,
		), nil
	}

	uID := uuid.New()
	r := rSM.Request{
		Id:       &uID,
		JITBuild: jitFunc,
	}

	resp := a.rS.Request(&r)
	if resp.Err != nil {
		return nil
	}

	session := &m.Session{}
	err := json.Unmarshal(resp.Resp, session)
	if err != nil {
		return nil
	}

	return session
}

func (a *apiUtil) TestSession(s *m.Session) string {
	jitFunc := func(args []interface{}) (string, error) {
		baseURL := "http://api.smitegame.com/smiteapi.svc/testsessionjson"
		t := time.Now().UTC()
		tS := a.aS.GetTimestamp(t)
		return fmt.Sprintf(
			"%s/%s/%s/%s/%s",
			baseURL,
			a.aS.GetID(),
			a.aS.GetSignature("testsession", tS),
			*s.SessionID,
			tS,
		), nil
	}

	uID := uuid.New()
	r := rSM.Request{
		Id:       &uID,
		JITBuild: jitFunc,
	}

	resp := a.rS.Request(&r)
	if resp.Err != nil {
		return ""
	}

	return string(resp.Resp)
}

func (a *apiUtil) GetDataUsed(s *m.Session) *m.UsageInfo {
	return nil
}
