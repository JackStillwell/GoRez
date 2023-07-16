package gorez

import (
	"fmt"
	"math"
	"strings"

	"github.com/JackStillwell/GoRez/internal"
	c "github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"

	authI "github.com/JackStillwell/GoRez/internal/auth/interfaces"

	requestI "github.com/JackStillwell/GoRez/internal/request/interfaces"
	requestM "github.com/JackStillwell/GoRez/internal/request/models"

	sessionI "github.com/JackStillwell/GoRez/internal/session/interfaces"
	sessionM "github.com/JackStillwell/GoRez/internal/session/models"
)

type matchInfo struct {
	authSvc authI.Service
	rqstSvc requestI.Service
	sesnSvc sessionI.Service
	gUtil   i.GorezUtil
	hrC     c.HiRezConstants
}

func NewMatchInfo(
	rS requestI.Service,
	aS authI.Service,
	sS sessionI.Service,
) i.MatchInfo {
	return &matchInfo{
		rqstSvc: rS,
		authSvc: aS,
		sesnSvc: sS,
		gUtil:   NewGorezUtil(aS, rS, sS),
		hrC:     c.NewHiRezConstants(),
	}
}

// GetMatchDetails will return data for players in a completed match
func (r *matchInfo) GetMatchDetails(matchID string) ([]byte, error) {
	mIds, errs := r.GetMatchDetailsBatch(matchID)
	return mIds[0], errs[0]
}

func (r *matchInfo) GetMatchDetailsBatch(matchIDs ...string) ([][]byte, []error) {
	requests := make([]func(*sessionM.Session) *requestM.Request, 0, (len(matchIDs)/10)+1)
	for i := len(matchIDs); i > 0; i = i - 10 {
		matchIdString := strings.Join(matchIDs[int(math.Max(0, float64(i-10))):i], ",")
		requestFunc := func(session *sessionM.Session) *requestM.Request {
			f := HiRezJIT(
				r.hrC.SmiteURLBase+"/"+r.hrC.GetMatchDetailsBatch+"json",
				r.authSvc.GetID(),
				r.hrC.GetMatchDetailsBatch,
				session.Key,
				r.authSvc.GetTimestamp,
				r.authSvc.GetSignature,
				matchIdString,
			)

			return &requestM.Request{JITFunc: f}
		}

		requests = append(requests, requestFunc)
	}

	// FIXME: need to check ret_msgs and assign errs accordingly

	return r.gUtil.BulkAsyncSessionRequest(requests)
}

// GetMatchDetails will return data for players in completed matches
func (r *matchInfo) GetMatchIDsByQueue(dateStrings []string, queueIDs []m.QueueID) (
	[]*[]m.MatchIDWithQueue, []error,
) {
	retObjs := []*[]m.MatchIDWithQueue{}
	errs := []error{}
	for i := 0; i < len(queueIDs); i++ {
		queueID := queueIDs[i]
		requests := make([]func(*sessionM.Session) *requestM.Request, 0,
			len(queueIDs)*len(dateStrings),
		)
		for j := 0; j < len(dateStrings); j++ {
			dateString := dateStrings[j]
			requestFunc := func(session *sessionM.Session) *requestM.Request {
				f := HiRezJIT(
					r.hrC.SmiteURLBase+"/"+r.hrC.GetMatchIDsByQueue+"json",
					r.authSvc.GetID(),
					r.hrC.GetMatchIDsByQueue,
					session.Key,
					r.authSvc.GetTimestamp,
					r.authSvc.GetSignature,
					fmt.Sprintf("%d", queueID)+"/"+dateString,
				)

				return &requestM.Request{JITFunc: f}
			}
			requests = append(requests, requestFunc)
		}
		rawObjs, requestErrs := r.gUtil.BulkAsyncSessionRequest(requests)
		unmarshaledObjs, unmarshalErrs := internal.UnmarshalObjs[[]m.MatchIDWithQueue](
			rawObjs, requestErrs,
		)
		for i := range unmarshaledObjs {
			matchIdsPtr := unmarshaledObjs[i]
			if matchIdsPtr == nil {
				continue
			}
			matchIds := *matchIdsPtr
			for j := range matchIds {
				matchIds[j].QueueID = int(queueID)
			}
		}
		retObjs = append(retObjs, unmarshaledObjs...)
		errs = append(errs, unmarshalErrs...)
	}

	return retObjs, errs
}

// GetMatchPlayerDetails will return data for players in a live match
func (r *matchInfo) GetMatchPlayerDetails(matchID string) ([]byte, error) {
	return r.GetMatchDetails(matchID)
}
