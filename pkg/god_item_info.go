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

type godItemInfo struct {
	authSvc authService.AuthService
	rqstSvc requestService.RequestService
	sesnSvc sessionService.SessionService
}

func NewGodItemInfo(
	rS requestService.RequestService,
	aS authService.AuthService,
	sS sessionService.SessionService,
) i.GodItemInfo {
	return &godItemInfo{
		rqstSvc: rS,
		authSvc: aS,
		sesnSvc: sS,
	}
}

func (g *godItemInfo) GetGods() ([]*m.God, error) {
	sessions, err := g.sesnSvc.ReserveSession(1)
	if err != nil {
		return nil, errors.Wrap(err, "reserving session")
	}

	defer g.sesnSvc.ReleaseSession(sessions)
	s := sessions[0]

	r := rSM.Request{
		JITArgs: []interface{}{
			hRConst.SmiteURLBase + hRConst.GetGods + "json",
			g.authSvc.GetID(),
			hRConst.GetGods,
			s.Key,
			g.authSvc.GetTimestamp,
			g.authSvc.GetSignature,
			"",
		},
		JITBuild: requestUtils.JITBase,
	}

	resp := g.rqstSvc.Request(&r)

	if resp.Err != nil {
		return nil, errors.Wrap(resp.Err, "requesting response")
	}

	gods := []*m.God{}
	err = json.Unmarshal(resp.Resp, gods)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling response")
	}

	return gods, nil
}

func (g *godItemInfo) GetItems() ([]*m.Item, error) {
	sessions, err := g.sesnSvc.ReserveSession(1)
	if err != nil {
		return nil, errors.Wrap(err, "reserving session")
	}

	defer g.sesnSvc.ReleaseSession(sessions)
	s := sessions[0]

	r := rSM.Request{
		JITArgs: []interface{}{
			hRConst.SmiteURLBase + hRConst.GetItems + "json",
			g.authSvc.GetID(),
			hRConst.GetItems,
			s.Key,
			g.authSvc.GetTimestamp,
			g.authSvc.GetSignature,
			"",
		},
		JITBuild: requestUtils.JITBase,
	}

	resp := g.rqstSvc.Request(&r)

	if resp.Err != nil {
		return nil, errors.Wrap(resp.Err, "requesting response")
	}

	items := []*m.Item{}
	err = json.Unmarshal(resp.Resp, items)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling response")
	}

	return items, nil
}

func (g *godItemInfo) GetGodRecItems(godIDs []int) ([]*m.ItemRecommendation, []error) {
	r := rSM.Request{
		JITArgs: []interface{}{
			hRConst.SmiteURLBase + hRConst.TestSession + "json",
			g.authSvc.GetID(),
			hRConst.TestSession,
			"",
			g.authSvc.GetTimestamp,
			g.authSvc.GetSignature,
			"",
		},
		JITBuild: requestUtils.JITBase,
	}

	uIDs := make([]*uuid.UUID, len(godIDs))
	for i := 0; i < len(godIDs); i++ {
		uID := uuid.New()
		r.Id = &uID
		g.rqstSvc.MakeRequest(&r)
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
