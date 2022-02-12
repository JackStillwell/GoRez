package gorez

import (
	"github.com/pkg/errors"

	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"

	authService "github.com/JackStillwell/GoRez/internal/auth_service/interfaces"
	requestService "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	sessionService "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
)

type matchInfo struct {
	authSvc authService.AuthService
	rqstSvc requestService.RequestService
	sesnSvc sessionService.SessionService
}

func NewMatchInfo(
	rS requestService.RequestService,
	aS authService.AuthService,
	sS sessionService.SessionService,
) i.MatchInfo {
	return &matchInfo{
		rqstSvc: rS,
		authSvc: aS,
		sesnSvc: sS,
	}
}

func (r *matchInfo) GetMatchDetails(matchID int) (*m.MatchDetails, error) {
	return nil, errors.New("unimplemented")
}

func (r *matchInfo) GetMatchDetailsBatch(matchIDs []int) ([]*m.MatchDetails, []error) {
	return nil, []error{errors.New("unimplemented")}
}

func (r *matchInfo) GetMatchIDsByQueue(queueID []m.QueueID) ([]*m.MatchIDWithQueue, []error) {
	return nil, []error{errors.New("unimplemented")}
}

func (r *matchInfo) GetMatchPlayerDetails(matchID []int) ([]*m.MatchDetails, []error) {
	return nil, []error{errors.New("unimplemented")}
}
