package session

import (
	"fmt"
	"sync"

	"github.com/JackStillwell/GoRez/internal/base"
	"go.uber.org/zap"

	i "github.com/JackStillwell/GoRez/internal/session/interfaces"
	m "github.com/JackStillwell/GoRez/internal/session/models"
)

type service struct {
	maxSessions       int
	availableSessions chan *m.Session
	reservedSessions  map[*m.Session]struct{}
	lock              sync.Mutex
	base              base.Service
}

func (s *service) addReservedSession(sess *m.Session) {
	log := s.base.GetLogger()
	log.Debug("waiting for reserved session lock", zap.String("operation", "addReservedSession"),
		zap.String("session", sess.Key))
	s.lock.Lock()
	log.Debug("reserved session lock acquired", zap.String("operation", "addReservedSession"),
		zap.String("session", sess.Key))
	s.reservedSessions[sess] = struct{}{}
	log.Debug("releasing reserved session lock", zap.String("operation", "addReservedSession"),
		zap.String("session", sess.Key))
	s.lock.Unlock()
}

func (s *service) removeReservedSession(sess *m.Session) {
	log := s.base.GetLogger()
	log.Debug("waiting for reserved session lock", zap.String("operation", "removeReservedSession"),
		zap.String("session", sess.Key))
	s.lock.Lock()
	log.Debug("reserved session lock acquired", zap.String("operation", "removeReservedSession"),
		zap.String("session", sess.Key))
	delete(s.reservedSessions, sess)
	log.Debug("releasing reserved session lock", zap.String("operation", "removeReservedSession"),
		zap.String("session", sess.Key))
	s.lock.Unlock()
}

func NewService(maxSessions int, existingSessions []*m.Session, b base.Service) i.Service {
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
		base:              b,
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
	log := s.base.GetLogger()
	log.Debug("reserving available session(s)", zap.Int("numReserving", numSessions),
		zap.Int("numAvailable", len(s.availableSessions)))
	for i := 0; i < numSessions; i++ {
		toReturn := <-s.availableSessions
		s.addReservedSession(toReturn)
		retChan <- toReturn
		log.Debug("reserved available session", zap.String("session", toReturn.Key),
			zap.Int("numReserved", numSessions), zap.Int("numAvailable", len(s.availableSessions)))
	}
}

func (s *service) ReleaseSession(sessions []*m.Session) {
	for _, session := range sessions {
		s.removeReservedSession(session)
	}

	for i := range sessions {
		s.availableSessions <- sessions[i]
	}
}

func (s *service) BadSession(sessions []*m.Session) {
	for _, session := range sessions {
		s.removeReservedSession(session)
	}
}
