package interfaces

import m "github.com/JackStillwell/GoRez/pkg/models"

type GoRez interface {
	GodItemInfo
	PlayerInfo
	MatchInfo
}

type APIUtil interface {
	CreateSession() *m.Session
	TestSession(*m.Session) string
	GetDataUsed(*m.Session) *m.UsageInfo
}

type GodItemInfo interface {
	GetGods() []*m.God
	GetItems() []*m.Item
	GetGodRecItems(godID int) []*m.ItemRecommendation
}

type PlayerInfo interface {
	GetPlayerIDByName(playerName string) *m.PlayerID
	GetPlayer(playerID int) *m.Player
	GetPlayerBatch(playerIDs []int) []*m.Player
	GetMatchHistory(playerID int) []*m.MatchDetails
	GetQueueStats(playerID int) []*m.QueueStat
}

type MatchInfo interface {
	GetMatchDetails(matchID int) *m.MatchDetails
	GetMatchDetailsBatch(matchIDs []int) []*m.MatchDetails
	GetMatchIDsByQueue(queueID m.QueueID) []*m.MatchID
	GetMatchPlayerDetails(matchID int) *m.MatchDetails
}
