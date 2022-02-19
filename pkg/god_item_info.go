package gorez

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/JackStillwell/GoRez/internal"

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

type godItemInfo struct {
	hrC c.HiRezConstants

	authSvc authI.AuthService
	rqstSvc requestI.RequestService
	sesnSvc sessionI.SessionService
}

func NewGodItemInfo(
	hrC c.HiRezConstants,
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

func (g *godItemInfo) singleRequest(r requestM.Request, unmarshalTo interface{}) error {
	sesnChan := make(chan *sessionM.Session, 1)
	g.sesnSvc.ReserveSession(1, sesnChan)
	s := <-sesnChan

	sessions := []*sessionM.Session{s}
	defer g.sesnSvc.ReleaseSession(sessions)

	r.JITArgs[1] = g.authSvc.GetID()
	r.JITArgs[3] = s.Key
	r.JITArgs[4] = g.authSvc.GetTimestamp
	r.JITArgs[5] = g.authSvc.GetSignature

	resp := g.rqstSvc.Request(&r)

	if resp.Err != nil {
		return errors.Wrap(resp.Err, "requesting response")
	}

	err := json.Unmarshal(resp.Resp, &unmarshalTo)
	if err != nil {
		return errors.Wrap(err, "marshaling response")
	}

	return nil
}

func (g *godItemInfo) GetGods() ([]*m.God, error) {
	r := requestM.Request{
		JITArgs: []interface{}{
			g.hrC.SmiteURLBase + g.hrC.GetGods + "json", "", g.hrC.GetGods, "", "", "", "",
		},
		JITBuild: requestU.JITBase,
	}

	gods := []*m.God{}
	err := g.singleRequest(r, &gods)
	return gods, err
}

func (g *godItemInfo) GetItems() ([]*m.Item, error) {
	r := requestM.Request{
		JITArgs: []interface{}{
			g.hrC.SmiteURLBase + g.hrC.GetItems + "json", "", g.hrC.GetItems, "", "", "", "",
		},
		JITBuild: requestU.JITBase,
	}

	items := []*m.Item{}
	err := g.singleRequest(r, &items)
	return items, err
}

func (g *godItemInfo) multiRequest(requestArgs []string, endpoint, method string,
) ([][]byte, []error) {
	requestBuilders := make([]func(*sessionM.Session) *requestM.Request, len(requestArgs))

	for i, arg := range requestArgs {
		requestBuilders[i] = func(s *sessionM.Session) *requestM.Request {
			return &requestM.Request{
				JITArgs: []interface{}{
					g.hrC.SmiteURLBase + endpoint,
					g.authSvc.GetID(),
					method,
					s.Key,
					g.authSvc.GetTimestamp,
					g.authSvc.GetSignature,
					arg,
				},
				JITBuild: requestU.JITBase,
			}
		}
	}

	rawObjs, errs := internal.BulkAsyncSessionRequest(g.rqstSvc, g.sesnSvc, requestBuilders)

	return rawObjs, errs
}

func (g *godItemInfo) GetGodRecItems(godIDs []int) ([]*m.ItemRecommendation, []error) {
	args := make([]string, len(godIDs))
	for i, gid := range godIDs {
		args[i] = fmt.Sprint(gid) + "/1"
	}

	endpoint := g.hrC.GetGodRecommendedItems + "json"
	rawObjs, errs := g.multiRequest(args, endpoint, g.hrC.GetGodRecommendedItems)

	itemRecs := make([]*m.ItemRecommendation, len(godIDs))
	for i, obj := range rawObjs {
		itemRec := &itemRecs[i]
		err := json.Unmarshal(obj, itemRec)
		if err != nil {
			errs[i] = errors.Wrap(err, "marshaling response")
		}
	}

	return itemRecs, errs
}
