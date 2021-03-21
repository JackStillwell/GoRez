package session_service

import (
	"fmt"
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

func (s *sessionService) ReserveSession(numSessions int) ([]*m.Session, error) {
	return nil, nil
}

func (s *sessionService) ReleaseSession(sessions []*m.Session) {

}

func (s *sessionService) BadSession(sessions []*m.Session) {
	remainingSessions := make([]*m.Session, 0, s.maxSessions)
outerLoop:
	for _, aS := range s.aSessions {
		for _, bS := range sessions {
			if aS.Key == bS.Key {
				continue outerLoop
			}
		}
		remainingSessions = append(remainingSessions, aS)
	}

	s.activeSessions = remainingSessions
}
