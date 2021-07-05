package session_service

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	i "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
	m "github.com/JackStillwell/GoRez/internal/session_service/models"
)

type sessionService struct {
	maxSessions       int
	availableSessions []*m.Session
	reservedSessions  []*m.Session
	lock              sync.Mutex
}

func NewSessionService(maxSessions int, existingSessions []*m.Session) (i.SessionService, error) {
	if len(existingSessions) > maxSessions {
		return nil, fmt.Errorf(
			"cannot create a session service with capacity %d and %d existing sessions",
			len(existingSessions),
			maxSessions,
		)
	}

	aS := make([]*m.Session, 0, maxSessions)
	copy(aS, existingSessions)

	rS := make([]*m.Session, 0, maxSessions)

	return &sessionService{
		maxSessions:       maxSessions,
		availableSessions: aS,
		reservedSessions:  rS,
		lock:              sync.Mutex{},
	}, nil
}

func (s *sessionService) ReserveSession(numSessions int, retChan chan *m.Session) {
	if numSessions > len(s.availableSessions) {
		// TODO: wait until available and return then
		return errors.New("not enough sessions available")
	}

	toReturn := s.availableSessions[:numSessions]
	s.availableSessions = s.availableSessions[numSessions-1:]
	s.reservedSessions = append(s.reservedSessions, toReturn...)

	return toReturn, nil
}

func (s *sessionService) ReleaseSession(sessions []*m.Session) {
	removeFromSlice(&s.reservedSessions, sessions)
	s.availableSessions = append(s.availableSessions, sessions...)
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

	sort.Ints(idxsToRemove)
	if idxsToRemove[len(idxsToRemove)-1] == len(*toModify)-1 {
		idxsToRemove = idxsToRemove[:len(idxsToRemove)-1]
		*toModify = (*toModify)[:len(idxsToRemove)-1]
	}
	for idx := len(idxsToRemove) - 1; idx >= 0; idx-- {
		iTR := idxsToRemove[idx]
		*toModify = append((*toModify)[:iTR], (*toModify)[iTR:]...)
	}
}
