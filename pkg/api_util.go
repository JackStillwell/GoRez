package gorez

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	c "github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"

	authI "github.com/JackStillwell/GoRez/internal/auth/interfaces"

	requestI "github.com/JackStillwell/GoRez/internal/request/interfaces"
	requestM "github.com/JackStillwell/GoRez/internal/request/models"

	sessionI "github.com/JackStillwell/GoRez/internal/session/interfaces"
	sessionM "github.com/JackStillwell/GoRez/internal/session/models"
)

type apiUtil struct {
	hiRezC  c.HiRezConstants
	authSvc authI.Service
	rqstSvc requestI.Service
	sesnSvc sessionI.Service
}

func NewAPIUtil(
	hrC c.HiRezConstants,
	aS authI.Service,
	rS requestI.Service,
	sS sessionI.Service,
) i.APIUtil {
	return &apiUtil{
		hiRezC:  hrC,
		authSvc: aS,
		rqstSvc: rS,
		sesnSvc: sS,
	}
}

func (a *apiUtil) CreateSession(numSessions int) ([]*m.Session, []error) {
	uIDs := make([]uuid.UUID, 0, numSessions)
	for i := 0; i < numSessions; i++ {
		r := requestM.Request{
			JITFunc: HiRezJIT(
				a.hiRezC.SmiteURLBase+"/"+a.hiRezC.CreateSession+"json",
				a.authSvc.GetID(),
				a.hiRezC.CreateSession,
				"",
				a.authSvc.GetTimestamp,
				a.authSvc.GetSignature,
				"",
			),
		}
		uID := uuid.New()
		r.Id = &uID
		a.rqstSvc.MakeRequest(&r)
		uIDs = append(uIDs, uID)
	}
	log.Println("create session requests made")
	log.Println("create session request uuids", uIDs)

	sessions := make([]*m.Session, 0, numSessions)
	errs := make([]error, 0, numSessions)
	for i := 0; i < numSessions; i++ {
		resp := a.rqstSvc.GetResponse(&uIDs[i])
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

func (a *apiUtil) TestSession(sessionKeys []string) ([]*string, []error) {
	uIDs := make([]uuid.UUID, 0, len(sessionKeys))
	for i := 0; i < len(sessionKeys); i++ {
		r := requestM.Request{
			JITFunc: HiRezJIT(
				a.hiRezC.SmiteURLBase+"/"+a.hiRezC.TestSession+"json",
				a.authSvc.GetID(),
				a.hiRezC.TestSession,
				sessionKeys[i],
				a.authSvc.GetTimestamp,
				a.authSvc.GetSignature,
				"",
			),
		}
		uID := uuid.New()
		r.Id = &uID
		a.rqstSvc.MakeRequest(&r)
		uIDs = append(uIDs, uID)
	}

	responses := make([]*string, 0, len(sessionKeys))
	errs := make([]error, 0, len(sessionKeys))
	for i := 0; i < len(sessionKeys); i++ {
		resp := a.rqstSvc.GetResponse(&uIDs[i])
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
		JITFunc: HiRezJIT(
			a.hiRezC.SmiteURLBase+"/"+a.hiRezC.GetDataUsed+"json",
			a.authSvc.GetID(),
			a.hiRezC.GetDataUsed,
			s.Key,
			a.authSvc.GetTimestamp,
			a.authSvc.GetSignature,
			"",
		),
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
