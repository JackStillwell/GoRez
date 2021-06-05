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
	GetGods() ([]*m.God, error)
	GetItems() ([]*m.Item, error)
	GetGodRecItems(godIDs []int) ([]*m.ItemRecommendation, []error)
}

type PlayerInfo interface {
	GetPlayerIDByName(playerName []string) *m.PlayerIDWithName
	GetPlayer(playerID int) *m.Player
	GetPlayerBatch(playerIDs []int) []*m.Player
	GetMatchHistory(playerID []int) []*m.MatchDetails
	GetQueueStats(playerID []int) []*m.QueueStat
}

type MatchInfo interface {
	GetMatchDetails(matchID int) *m.MatchDetails
	GetMatchDetailsBatch(matchIDs []int) []*m.MatchDetails
	GetMatchIDsByQueue(queueID []m.QueueID) []*m.MatchIDWithQueue
	GetMatchPlayerDetails(matchID []int) []*m.MatchDetails
}
