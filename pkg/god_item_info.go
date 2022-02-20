package gorez

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	c "github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"

	requestM "github.com/JackStillwell/GoRez/internal/request_service/models"
	requestU "github.com/JackStillwell/GoRez/internal/request_service/utilities"
)

type godItemInfo struct {
	hrC c.HiRezConstants

	util i.GorezUtil
}

func NewGodItemInfo(
	hrC c.HiRezConstants,
	uG i.GorezUtil,
) i.GodItemInfo {
	return &godItemInfo{
		hrC:  hrC,
		util: uG,
	}
}

func (g *godItemInfo) GetGods() ([]*m.God, error) {
	r := requestM.Request{
		JITArgs: []interface{}{
			g.hrC.SmiteURLBase + g.hrC.GetGods + "json", "", g.hrC.GetGods, "", "", "", "",
		},
		JITBuild: requestU.JITBase,
	}

	gods := []*m.God{}
	err := g.util.SingleRequest(r, &gods)
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
	err := g.util.SingleRequest(r, &items)
	return items, err
}

func (g *godItemInfo) GetGodRecItems(godIDs []int) ([]*m.ItemRecommendation, []error) {
	args := make([]string, len(godIDs))
	for i, gid := range godIDs {
		args[i] = fmt.Sprint(gid) + "/1"
	}

	baseURL := g.hrC.SmiteURLBase + "/" + g.hrC.GetGodRecommendedItems + "json"
	rawObjs, errs := g.util.MultiRequest(args, baseURL, g.hrC.GetGodRecommendedItems)

	itemRecs := make([]*m.ItemRecommendation, len(godIDs))
	for i, obj := range rawObjs {
		if obj != nil {
			itemRec := itemRecs[i]
			err := json.Unmarshal(obj, itemRec)
			if err != nil {
				errs[i] = errors.Wrap(err, "unmarshaling response")
			}
		}
	}

	return itemRecs, errs
}
