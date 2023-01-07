package gorez

import (
	"encoding/json"
	"log"

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
	hiRezC  c.HiRezConstants
	authSvc authI.AuthService
	rqstSvc requestI.RequestService
	sesnSvc sessionI.SessionService
}

func NewAPIUtil(
	hrC c.HiRezConstants,
	aS authI.AuthService,
	rS requestI.RequestService,
	sS sessionI.SessionService,
) i.APIUtil {
	return &apiUtil{
		hiRezC:  hrC,
		authSvc: aS,
		rqstSvc: rS,
		sesnSvc: sS,
	}
}

func (a *apiUtil) CreateSession(numSessions int) ([]*m.Session, []error) {
	r := requestM.Request{
		JITArgs: []any{
			a.hiRezC.SmiteURLBase + "/" + a.hiRezC.CreateSession + "json",
			a.authSvc.GetID(),
			a.hiRezC.CreateSession,
			"",
			a.authSvc.GetTimestamp,
			a.authSvc.GetSignature,
			"",
		},
		JITBuild: requestU.JITBase,
	}

	uIDs := make([]*uuid.UUID, 0, numSessions)
	for i := 0; i < numSessions; i++ {
		uID := uuid.New()
		r.Id = &uID
		a.rqstSvc.MakeRequest(&r)
		uIDs = append(uIDs, &uID)
	}
	log.Println("create session requests made")
	log.Println("create session request uuids", uIDs)

	log.Println("getting create session responses")
	responseChan := make(chan *requestM.RequestResponse, numSessions)
	for i := 0; i < numSessions; i++ {
		a.rqstSvc.GetResponse(uIDs[i], responseChan)
	}
	log.Println("create session responses received")

	sessions := make([]*m.Session, 0, numSessions)
	errs := make([]error, 0, numSessions)
	for i := 0; i < numSessions; i++ {
		resp := <-responseChan
		if resp.Err != nil {
			sessions = append(sessions, nil)
			errs = append(errs, errors.Wrap(resp.Err, "request"))
			continue
		}

		session := &m.Session{}
		err := json.Unmarshal(resp.Resp, session)
		if err != nil {
			sessions = append(sessions, nil)
			errs = append(errs, errors.Wrap(err, "unmarshal response"))
			continue
		}

		sessions = append(sessions, session)
		errs = append(errs, nil)
	}

	return sessions, errs
}

func (a *apiUtil) TestSession(s []*m.Session) ([]*string, []error) {
	r := requestM.Request{
		JITArgs: []interface{}{
			a.hiRezC.SmiteURLBase + "/" + a.hiRezC.TestSession + "json",
			a.authSvc.GetID(),
			a.hiRezC.TestSession,
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
			a.hiRezC.SmiteURLBase + "/" + a.hiRezC.GetDataUsed + "json",
			a.authSvc.GetID(),
			a.hiRezC.GetDataUsed,
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
