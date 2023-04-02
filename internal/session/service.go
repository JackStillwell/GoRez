package session

import (
	"fmt"
	"log"
	"sort"
	"sync"

	i "github.com/JackStillwell/GoRez/internal/session/interfaces"
	m "github.com/JackStillwell/GoRez/internal/session/models"
)

type service struct {
	maxSessions       int
	availableSessions chan *m.Session
	reservedSessions  []*m.Session
	lock              sync.Mutex
}

func NewService(maxSessions int, existingSessions []*m.Session) i.Service {
	if len(existingSessions) > maxSessions {
		panic(fmt.Sprintf(
			"cannot create a session service with capacity %d and %d existing sessions",
			maxSessions,
			len(existingSessions),
		))
	}

	aS := make(chan *m.Session, maxSessions)
	for i := range existingSessions {
		aS <- existingSessions[i]
	}

	rS := make([]*m.Session, 0, maxSessions)

	return &service{
		maxSessions:       maxSessions,
		availableSessions: aS,
		reservedSessions:  rS,
		lock:              sync.Mutex{},
	}
}

func (s *service) GetAvailableSessions() []*m.Session {
	numSessions := len(s.availableSessions)
	toRet := make([]*m.Session, 0, numSessions)
	for i := 0; i < numSessions; i++ {
		toRet = append(toRet, <-s.availableSessions)
	}
	return toRet
}

func (s *service) ReserveSession(numSessions int, retChan chan *m.Session) {
	log.Printf("reserving %d of %d available sessions", numSessions, len(s.availableSessions))
	for i := 0; i < numSessions; i++ {
		toReturn := <-s.availableSessions
		s.reservedSessions = append(s.reservedSessions, toReturn)
		retChan <- toReturn
		log.Printf("%d of %d sessions reserved", i+1, numSessions)
	}
}

func (s *service) ReleaseSession(sessions []*m.Session) {
	removeFromSlice(&s.reservedSessions, sessions)

	for i := range sessions {
		s.availableSessions <- sessions[i]
	}
}

func (s *service) BadSession(sessions []*m.Session) {
	removeFromSlice(&s.reservedSessions, sessions)
}

func removeFromSlice(toModify *[]*m.Session, toRemove []*m.Session) {
	idxsToRemove := make([]int, 0, len(toRemove))
	for idxRs, rS := range *toModify {
		for _, bS := range toRemove {
			if rS.Key == bS.Key {
				idxsToRemove = append(idxsToRemove, idxRs)
			}
		}
	}

	// None of the sessions to be removed were present in the slice to be modified
	if len(idxsToRemove) == 0 {
		return
	}

	sort.Ints(idxsToRemove)
	if idxsToRemove[len(idxsToRemove)-1] == len(*toModify)-1 {
		*toModify = (*toModify)[:len(idxsToRemove)-1]
		idxsToRemove = idxsToRemove[:len(idxsToRemove)-1]
	}
	for idx := len(idxsToRemove) - 1; idx >= 0; idx-- {
		iTR := idxsToRemove[idx]
		*toModify = append((*toModify)[:iTR], (*toModify)[iTR:]...)
	}
}
