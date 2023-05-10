package session

import (
	"fmt"
	"log"
	"sync"

	i "github.com/JackStillwell/GoRez/internal/session/interfaces"
	m "github.com/JackStillwell/GoRez/internal/session/models"
)

type service struct {
	maxSessions       int
	availableSessions chan *m.Session
	reservedSessions  map[*m.Session]struct{}
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

	rS := make(map[*m.Session]struct{}, maxSessions)

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
		s.lock.Lock()
		s.reservedSessions[toReturn] = struct{}{}
		s.lock.Unlock()
		retChan <- toReturn
		log.Printf("%d of %d sessions reserved", i+1, numSessions)
	}
}

func (s *service) ReleaseSession(sessions []*m.Session) {
	s.lock.Lock()
	for _, session := range sessions {
		delete(s.reservedSessions, session)
	}
	s.lock.Unlock()

	for i := range sessions {
		s.availableSessions <- sessions[i]
	}
}

func (s *service) BadSession(sessions []*m.Session) {
	s.lock.Lock()
	for _, session := range sessions {
		delete(s.reservedSessions, session)
	}
	s.lock.Unlock()
}
