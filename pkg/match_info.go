package gorez

import (
	"github.com/pkg/errors"

	c "github.com/JackStillwell/GoRez/pkg/constants"
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
	gUtil   i.GorezUtil
	hrC     c.HiRezConstants
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
		gUtil:   NewGorezUtil(aS, rS, sS),
		hrC:     c.NewHiRezConstants(),
	}
}

func (r *matchInfo) GetMatchDetails(matchID int) (*m.MatchDetails, error) {
	return nil, errors.New("unimplemented")
}

func (r *matchInfo) GetMatchDetailsBatch(matchIDs []int) ([]*m.MatchDetails, []error) {

	/*requests := make([]requestM.Request, 0, (len(matchIDs)/10)+1)
	for i := len(matchIDs); i > 0; i = i - 10 {
		builder, err := requestU.JITBase(
			r.hrC.SmiteURLBase+"/"+r.hrC.GetMatchDetailsBatch+"json",
			r.authSvc.GetID(),
			r.hrC.GetMatchDetailsBatch,
			"",
			"",
			"",
			"",
		)
		requests = append(requests, builder)
	}*/

	return nil, []error{errors.New("unimplemented")}
}

func (r *matchInfo) GetMatchIDsByQueue(queueID []m.QueueID) ([]*m.MatchIDWithQueue, []error) {
	return nil, []error{errors.New("unimplemented")}
}

func (r *matchInfo) GetMatchPlayerDetails(matchID []int) ([]*m.MatchDetails, []error) {
	return nil, []error{errors.New("unimplemented")}
}
