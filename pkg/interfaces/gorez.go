package interfaces

import m "github.com/JackStillwell/GoRez/pkg/models"

//go:generate mockgen --source=gorez.go --destination=../mocks/mock_gorez.go --package=mock

type GoRez interface {
	Init(numSessions int) error
	Shutdown()
	GodItemInfo
	PlayerInfo
	MatchInfo
}

type APIUtil interface {
	CreateSession(int) ([]*m.Session, []error)
	TestSession([]string) ([]*string, []error)
	GetDataUsed() (*m.UsageInfo, error)
}

type GodItemInfo interface {
	GetGods() ([]byte, error)
	GetItems() ([]byte, error)
	GetGodRecItems(godIDs []int) ([][]byte, []error)
}

type PlayerInfo interface {
	GetPlayerIDByName(playerName []string) ([]*m.PlayerIDWithName, error)
	GetPlayer(playerID int) ([]byte, error)
	GetPlayerBatch(playerIDs []int) ([][]byte, []error)
	GetMatchHistory(playerID []int) ([][]byte, []error)
	GetQueueStats(playerID []int) ([]*m.QueueStat, []error)
}

type MatchInfo interface {
	GetMatchDetails(matchID string) ([]byte, error)
	GetMatchDetailsBatch(matchIDs ...string) ([][]byte, []error)
	GetMatchIDsByQueue(dateStrings []string, queueIDs []m.QueueID) ([]*[]m.MatchIDWithQueue, []error)
	GetMatchPlayerDetails(matchID string) ([]byte, error)
}
