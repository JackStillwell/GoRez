package gorez

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/JackStillwell/GoRez/internal"
	authService "github.com/JackStillwell/GoRez/internal/auth_service/interfaces"
	requestService "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	rSM "github.com/JackStillwell/GoRez/internal/request_service/models"
	requestUtils "github.com/JackStillwell/GoRez/internal/request_service/utilities"
	sessionService "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
	sSM "github.com/JackStillwell/GoRez/internal/session_service/models"
	hRConst "github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"
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
	sesnChan := make(chan *sSM.Session, 1)
	g.sesnSvc.ReserveSession(1, sesnChan)
	s := <-sesnChan

	sessions := []*sSM.Session{s}
	defer g.sesnSvc.ReleaseSession(sessions)

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
	err := json.Unmarshal(resp.Resp, &gods)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling response")
	}

	return gods, nil
}

func (g *godItemInfo) GetItems() ([]*m.Item, error) {
	sesnChan := make(chan *sSM.Session, 1)
	g.sesnSvc.ReserveSession(1, sesnChan)
	s := <-sesnChan

	sessions := []*sSM.Session{s}
	defer g.sesnSvc.ReleaseSession(sessions)

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
	err := json.Unmarshal(resp.Resp, &items)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling response")
	}

	return items, nil
}

func (g *godItemInfo) GetGodRecItems(godIDs []int) ([]*m.ItemRecommendation, []error) {
	requestBuilders := make([]func(*sSM.Session) *rSM.Request, len(godIDs))

	for i, gid := range godIDs {
		requestBuilders[i] = func(s *sSM.Session) *rSM.Request {
			return &rSM.Request{
				JITArgs: []interface{}{
					hRConst.SmiteURLBase + hRConst.GetGodRecommendedItems + "json",
					g.authSvc.GetID(),
					hRConst.GetGodRecommendedItems,
					s.Key,
					g.authSvc.GetTimestamp,
					g.authSvc.GetSignature,
					fmt.Sprint(gid) + "/1",
				},
				JITBuild: requestUtils.JITBase,
			}
		}
	}

	rawObjs, errs := internal.BulkAsyncSessionRequest(g.rqstSvc, g.sesnSvc, requestBuilders,
		func(b []byte) (interface{}, error) {
			itemRec := &m.ItemRecommendation{}
			err := json.Unmarshal(b, itemRec)
			if err != nil {
				return nil, errors.Wrap(err, "marshaling response")
			}
			return itemRec, nil
		})

	itemRecs := make([]*m.ItemRecommendation, len(godIDs))
	for i, obj := range rawObjs {
		rec, ok := obj.(*m.ItemRecommendation)
		if !ok {
			errs = append(errs, errors.New("converting from interface to itemrecommendation"))
		}

		itemRecs[i] = rec
	}

	return itemRecs, errs
}
