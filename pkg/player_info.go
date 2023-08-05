package gorez

import (
	"fmt"

	"github.com/JackStillwell/GoRez/internal"
	"github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"

	authI "github.com/JackStillwell/GoRez/internal/auth/interfaces"
	requestI "github.com/JackStillwell/GoRez/internal/request/interfaces"
	requestM "github.com/JackStillwell/GoRez/internal/request/models"
	sessionI "github.com/JackStillwell/GoRez/internal/session/interfaces"
	sessionM "github.com/JackStillwell/GoRez/internal/session/models"
)

type playerInfo struct {
	authSvc authI.Service
	rqstSvc requestI.Service
	sesnSvc sessionI.Service
	gUtil   i.GorezUtil
	hrC     constants.HiRezConstants
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
		gUtil:   NewGorezUtil(aS, rS, sS),
		hrC:     constants.NewHiRezConstants(),
	}
}

func (r *playerInfo) GetPlayerIDByName(playerNames ...string) ([]*m.PlayerIDWithName, []error) {
	requests := make([]func(*sessionM.Session) *requestM.Request, 0, len(playerNames))
	for i := range playerNames {
		method := r.hrC.GetPlayerIDByName
		requestFunc := func(session *sessionM.Session) *requestM.Request {
			f := HiRezJIT(
				r.hrC.SmiteURLBase+"/"+method+"json",
				r.authSvc.GetID(),
				method,
				session.Key,
				r.authSvc.GetTimestamp,
				r.authSvc.GetSignature,
				playerNames[i],
			)

			return &requestM.Request{JITFunc: f}
		}

		requests = append(requests, requestFunc)
	}

	// FIXME: need to check ret_msgs and assign errs accordingly
	return internal.UnmarshalObjs[m.PlayerIDWithName](r.gUtil.BulkAsyncSessionRequest(requests))
}

func (r *playerInfo) GetPlayer(playerID int) ([]byte, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (r *playerInfo) GetPlayerBatch(playerIDs []int) ([][]byte, []error) {
	return nil, []error{fmt.Errorf("unimplemented")}
}

func (r *playerInfo) GetMatchHistory(playerIDs ...int) ([][]byte, []error) {
	requests := make([]func(*sessionM.Session) *requestM.Request, 0, len(playerIDs))
	for i := range playerIDs {
		method := r.hrC.GetMatchHistory
		requestFunc := func(session *sessionM.Session) *requestM.Request {
			f := HiRezJIT(
				r.hrC.SmiteURLBase+"/"+method+"json",
				r.authSvc.GetID(),
				method,
				session.Key,
				r.authSvc.GetTimestamp,
				r.authSvc.GetSignature,
				fmt.Sprintf("%d", playerIDs[i]),
			)

			return &requestM.Request{JITFunc: f}
		}

		requests = append(requests, requestFunc)
	}

	// FIXME: need to check ret_msgs and assign errs accordingly

	return r.gUtil.BulkAsyncSessionRequest(requests)
}

func (r *playerInfo) GetQueueStats(playerID []int) ([]*m.QueueStat, []error) {
	return nil, []error{fmt.Errorf("unimplemented")}
}
