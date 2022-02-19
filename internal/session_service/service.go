package session_service

import (
	"fmt"
	"sort"
	"sync"

	i "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
	m "github.com/JackStillwell/GoRez/internal/session_service/models"
)

type sessionService struct {
	maxSessions       int
	availableSessions chan *m.Session
	reservedSessions  []*m.Session
	lock              sync.Mutex
}

func NewSessionService(maxSessions int, existingSessions []*m.Session) i.SessionService {
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

	return &sessionService{
		maxSessions:       maxSessions,
		availableSessions: aS,
		reservedSessions:  rS,
		lock:              sync.Mutex{},
	}
}

func (s *sessionService) GetAvailableSessions() []*m.Session {
	numSessions := len(s.availableSessions)
	toRet := make([]*m.Session, 0, numSessions)
	for i := 0; i < numSessions; i++ {
		toRet = append(toRet, <-s.availableSessions)
	}
	return toRet
}

func (s *sessionService) ReserveSession(numSessions int, retChan chan *m.Session) {
	for i := 0; i < numSessions; i++ {
		toReturn := <-s.availableSessions
		s.reservedSessions = append(s.reservedSessions, toReturn)
		retChan <- toReturn
	}
}

func (s *sessionService) ReleaseSession(sessions []*m.Session) {
	removeFromSlice(&s.reservedSessions, sessions)

	for i := range sessions {
		s.availableSessions <- sessions[i]
	}
}

func (s *sessionService) BadSession(sessions []*m.Session) {
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
