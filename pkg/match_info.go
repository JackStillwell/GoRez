package gorez

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/JackStillwell/GoRez/internal"
	c "github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"

	authService "github.com/JackStillwell/GoRez/internal/auth_service/interfaces"
	requestService "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	requestM "github.com/JackStillwell/GoRez/internal/request_service/models"
	requestU "github.com/JackStillwell/GoRez/internal/request_service/utilities"
	sessionService "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
	sessionM "github.com/JackStillwell/GoRez/internal/session_service/models"
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

// GetMatchDetails will return data for players in a completed match
func (r *matchInfo) GetMatchDetails(matchID int) (*m.MatchDetails, error) {
	mIds, errs := r.GetMatchDetailsBatch(matchID)
	return mIds[0], errs[0]
}

// GetMatchDetails will return data for players in completed matches
func (r *matchInfo) GetMatchDetailsBatch(matchIDs ...int) ([]*m.MatchDetails, []error) {
	requests := make([]func(*sessionM.Session) *requestM.Request, 0, (len(matchIDs)/10)+1)
	for i := len(matchIDs); i > 0; i = i - 10 {
		matchIdStrings := make([]string, 0, 10)
		for _, v := range matchIDs {
			matchIdStrings = append(matchIdStrings, fmt.Sprintf("%d", v))
		}
		requestFunc := func(session *sessionM.Session) *requestM.Request {
			args := []any{
				r.hrC.SmiteURLBase + "/" + r.hrC.GetMatchDetailsBatch + "json",
				r.authSvc.GetID(),
				r.hrC.GetMatchDetailsBatch,
				session.Key,
				r.authSvc.GetTimestamp,
				r.authSvc.GetSignature,
				strings.Join(matchIdStrings, ","),
			}

			return &requestM.Request{JITArgs: args, JITBuild: requestU.JITBase}
		}

		requests = append(requests, requestFunc)
	}

	rawObjs, errs := r.gUtil.BulkAsyncSessionRequest(requests)
	log.Println("rawObjs retrieved:", rawObjs)
	fmt.Println(errs)
	f, err := os.Create("rawobjs.json")
	if err != nil {
		log.Println("error opening file to write matchdetails", err)
	}
	f.Close()
	f.Write(rawObjs[0])

	nestedMDs, errs := internal.UnmarshalObjs[[]m.MatchDetails](rawObjs, errs)

	flatMDs := make([]*m.MatchDetails, 0, 10*len(nestedMDs))
	for _, v := range nestedMDs {
		for _, w := range *v {
			flatMDs = append(flatMDs, &w)
		}
	}

	return flatMDs, errs
}

func (r *matchInfo) GetMatchIDsByQueue(queueID []m.QueueID) ([]*m.MatchIDWithQueue, []error) {
	return nil, []error{errors.New("unimplemented")}
}

// GetMatchPlayerDetails will return data for players in a live match
func (r *matchInfo) GetMatchPlayerDetails(matchID int) (*m.MatchDetails, error) {
	return r.GetMatchDetails(matchID)
}
