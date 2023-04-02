package gorez

import (
	"fmt"

	c "github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
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

func (g *godItemInfo) GetGods() ([]byte, error) {
	return g.util.SingleRequest(g.hrC.SmiteURLBase+"/"+g.hrC.GetGods+"json", g.hrC.GetGods, "1")
}

func (g *godItemInfo) GetItems() ([]byte, error) {
	return g.util.SingleRequest(g.hrC.SmiteURLBase+"/"+g.hrC.GetItems+"json", g.hrC.GetItems, "1")
}

func (g *godItemInfo) GetGodRecItems(godIDs []int) ([][]byte, []error) {
	args := make([]string, len(godIDs))
	for i, gid := range godIDs {
		args[i] = fmt.Sprint(gid) + "/1"
	}

	baseURL := g.hrC.SmiteURLBase + "/" + g.hrC.GetGodRecommendedItems + "json"
	return g.util.MultiRequest(args, baseURL, g.hrC.GetGodRecommendedItems)
}
