package interfaces

import m "github.com/JackStillwell/GoRez/pkg/models"

//go:generate mockgen --source=gorez.go --destination=../mocks/mock_gorez.go --package=mock

type GoRez interface {
	Init() error
	Shutdown()
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
	GetPlayerIDByName(playerName []string) (*m.PlayerIDWithName, error)
	GetPlayer(playerID int) (*m.Player, error)
	GetPlayerBatch(playerIDs []int) ([]*m.Player, []error)
	GetMatchHistory(playerID []int) ([]*m.MatchDetails, []error)
	GetQueueStats(playerID []int) ([]*m.QueueStat, []error)
}

type MatchInfo interface {
	GetMatchDetails(matchID int) (*[]m.MatchDetails, error)
	GetMatchDetailsBatch(matchIDs ...int) ([]*[]m.MatchDetails, []error)
	GetMatchDetailsBatchRaw(matchIDs ...int) ([][]byte, []error)
	GetMatchIDsByQueue(dateStrings []string, queueIDs []m.QueueID) ([]*[]m.MatchIDWithQueue, []error)
	GetMatchIDsByQueueRaw(dateStrings []string, queueIDs []m.QueueID) ([][]byte, []error)
	GetMatchPlayerDetails(matchID int) (*[]m.MatchDetails, error)
}
