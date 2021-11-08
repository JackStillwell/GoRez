package gorez

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

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
	numIDs := len(godIDs)
	responseChan := make(chan *rSM.RequestResponse, numIDs)
	uIDSessionMap := make(map[*uuid.UUID]*sSM.Session, numIDs)

	// NOTE: this is async so the reservation and release of sessions is possible, but the func
	// return depends upon responses being completed.
	go func() {
		for _, gid := range godIDs {
			sessChan := make(chan *sSM.Session, 1)
			g.sesnSvc.ReserveSession(1, sessChan)
			s := <-sessChan // NOTE: will wait here until session recieved
			r := rSM.Request{
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

			uID := uuid.New()
			r.Id = &uID
			uIDSessionMap[&uID] = s
			g.rqstSvc.MakeRequest(&r)
			g.rqstSvc.GetResponse(&uID, responseChan)
		}
	}()

	responses := make([]*m.ItemRecommendation, 0, numIDs)
	errs := make([]error, 0, numIDs)
	for i := 0; i < numIDs; i++ {
		resp := <-responseChan
		if resp.Err != nil {
			if strings.Contains(resp.Err.Error(), "session") {
				g.sesnSvc.BadSession([]*sSM.Session{uIDSessionMap[resp.Id]})
			} else {
				g.sesnSvc.ReleaseSession([]*sSM.Session{uIDSessionMap[resp.Id]})
			}
			errs = append(errs, errors.Wrap(resp.Err, "request"))
			continue
		}

		g.sesnSvc.ReleaseSession([]*sSM.Session{uIDSessionMap[resp.Id]})

		itemRec := &m.ItemRecommendation{}
		err := json.Unmarshal(resp.Resp, itemRec)
		if err != nil {
			errs = append(errs, errors.Wrap(resp.Err, "marshaling response"))
			continue
		}
		responses = append(responses, itemRec)
	}

	return responses, errs
}
