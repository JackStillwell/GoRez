package interfaces

import m "github.com/JackStillwell/GoRez/pkg/models"

//go:generate mockgen --source=gorez.go --destination=../mocks/mock_gorez.go --package=mock

type GoRez interface {
	GodItemInfo
	PlayerInfo
	MatchInfo
}

type APIUtil interface {
	CreateSession(int) ([]*m.Session, []error)
	TestSession([]*m.Session) ([]*string, []error)
	GetDataUsed() (*m.UsageInfo, error)
}

type GodItemInfo interface {
	GetGods() []*m.God
	GetItems() []*m.Item
	// TODO: need a parent obj to include godID
	GetGodRecItems(godID []int) []*m.ItemRecommendation
}

type PlayerInfo interface {
	// TODO: need a parent obj to include playerName
	GetPlayerIDByName(playerName []string) *m.PlayerID
	GetPlayer(playerID int) *m.Player
	GetPlayerBatch(playerIDs []int) []*m.Player
	// TODO: need a parent obj to include playerID
	GetMatchHistory(playerID []int) []*m.MatchDetails
	// TODO: need a parent obj to include playerID
	GetQueueStats(playerID []int) []*m.QueueStat
}

type MatchInfo interface {
	GetMatchDetails(matchID int) *m.MatchDetails
	GetMatchDetailsBatch(matchIDs []int) []*m.MatchDetails
	// TODO: need a parent obj to include queueid
	GetMatchIDsByQueue(queueID []m.QueueID) []*m.MatchID
	GetMatchPlayerDetails(matchID []int) []*m.MatchDetails
}
