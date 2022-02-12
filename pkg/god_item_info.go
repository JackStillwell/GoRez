package gorez

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/JackStillwell/GoRez/internal"

	hRConst "github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"

	authI "github.com/JackStillwell/GoRez/internal/auth_service/interfaces"

	requestI "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	requestM "github.com/JackStillwell/GoRez/internal/request_service/models"
	requestU "github.com/JackStillwell/GoRez/internal/request_service/utilities"

	sessionI "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
	sessionM "github.com/JackStillwell/GoRez/internal/session_service/models"
)

type godItemInfo struct {
	hrC hRConst.HiRezConstants

	authSvc authI.AuthService
	rqstSvc requestI.RequestService
	sesnSvc sessionI.SessionService
}

func NewGodItemInfo(
	hrC hRConst.HiRezConstants,
	rS requestI.RequestService,
	aS authI.AuthService,
	sS sessionI.SessionService,
) i.GodItemInfo {
	return &godItemInfo{
		hrC: hrC,

		rqstSvc: rS,
		authSvc: aS,
		sesnSvc: sS,
	}
}

func (g *godItemInfo) GetGods() ([]*m.God, error) {
	sesnChan := make(chan *sessionM.Session, 1)
	g.sesnSvc.ReserveSession(1, sesnChan)
	s := <-sesnChan

	sessions := []*sessionM.Session{s}
	defer g.sesnSvc.ReleaseSession(sessions)

	r := requestM.Request{
		JITArgs: []interface{}{
			g.hrC.SmiteURLBase + g.hrC.GetGods + "json",
			g.authSvc.GetID(),
			g.hrC.GetGods,
			s.Key,
			g.authSvc.GetTimestamp,
			g.authSvc.GetSignature,
			"",
		},
		JITBuild: requestU.JITBase,
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
	sesnChan := make(chan *sessionM.Session, 1)
	g.sesnSvc.ReserveSession(1, sesnChan)
	s := <-sesnChan

	sessions := []*sessionM.Session{s}
	defer g.sesnSvc.ReleaseSession(sessions)

	r := requestM.Request{
		JITArgs: []interface{}{
			g.hrC.SmiteURLBase + g.hrC.GetItems + "json",
			g.authSvc.GetID(),
			g.hrC.GetItems,
			s.Key,
			g.authSvc.GetTimestamp,
			g.authSvc.GetSignature,
			"",
		},
		JITBuild: requestU.JITBase,
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
	requestBuilders := make([]func(*sessionM.Session) *requestM.Request, len(godIDs))

	for i, gid := range godIDs {
		requestBuilders[i] = func(s *sessionM.Session) *requestM.Request {
			return &requestM.Request{
				JITArgs: []interface{}{
					g.hrC.SmiteURLBase + g.hrC.GetGodRecommendedItems + "json",
					g.authSvc.GetID(),
					g.hrC.GetGodRecommendedItems,
					s.Key,
					g.authSvc.GetTimestamp,
					g.authSvc.GetSignature,
					fmt.Sprint(gid) + "/1",
				},
				JITBuild: requestU.JITBase,
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
