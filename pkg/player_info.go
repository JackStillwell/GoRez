package gorez

import (
	"github.com/pkg/errors"

	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"

	authService "github.com/JackStillwell/GoRez/internal/auth_service/interfaces"
	requestService "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	sessionService "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
)

type playerInfo struct {
	authSvc authService.AuthService
	rqstSvc requestService.RequestService
	sesnSvc sessionService.SessionService
}

func NewPlayerInfo(
	rS requestService.RequestService,
	aS authService.AuthService,
	sS sessionService.SessionService,
) i.PlayerInfo {
	return &playerInfo{
		rqstSvc: rS,
		authSvc: aS,
		sesnSvc: sS,
	}
}

func (r *playerInfo) GetPlayerIDByName(playerName []string) (*m.PlayerIDWithName, error) {
	return nil, errors.New("unimplemented")
}

func (r *playerInfo) GetPlayer(playerID int) (*m.Player, error) {
	return nil, errors.New("unimplemented")
}

func (r *playerInfo) GetPlayerBatch(playerIDs []int) ([]*m.Player, []error) {
	return nil, []error{errors.New("unimplemented")}
}

func (r *playerInfo) GetMatchHistory(playerID []int) ([]*m.MatchDetails, []error) {
	return nil, []error{errors.New("unimplemented")}
}

func (r *playerInfo) GetQueueStats(playerID []int) ([]*m.QueueStat, []error) {
	return nil, []error{errors.New("unimplemented")}
}
