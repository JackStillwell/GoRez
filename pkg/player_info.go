package gorez

import (
	"github.com/pkg/errors"

	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"

	authI "github.com/JackStillwell/GoRez/internal/auth/interfaces"
	requestI "github.com/JackStillwell/GoRez/internal/request/interfaces"
	sessionI "github.com/JackStillwell/GoRez/internal/session/interfaces"
)

type playerInfo struct {
	authSvc authI.Service
	rqstSvc requestI.Service
	sesnSvc sessionI.Service
}

func NewPlayerInfo(
	rS requestI.Service,
	aS authI.Service,
	sS sessionI.Service,
) i.PlayerInfo {
	return &playerInfo{
		rqstSvc: rS,
		authSvc: aS,
		sesnSvc: sS,
	}
}

func (r *playerInfo) GetPlayerIDByName(playerName []string) ([]*m.PlayerIDWithName, error) {
	return nil, errors.New("unimplemented")
}

func (r *playerInfo) GetPlayer(playerID int) ([]byte, error) {
	return nil, errors.New("unimplemented")
}

func (r *playerInfo) GetPlayerBatch(playerIDs []int) ([][]byte, []error) {
	return nil, []error{errors.New("unimplemented")}
}

func (r *playerInfo) GetMatchHistory(playerID []int) ([][]byte, []error) {
	return nil, []error{errors.New("unimplemented")}
}

func (r *playerInfo) GetQueueStats(playerID []int) ([]*m.QueueStat, []error) {
	return nil, []error{errors.New("unimplemented")}
}
