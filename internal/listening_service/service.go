package listening_service

import (
	"sync"

	i "github.com/JackStillwell/GoRez/internal/listening_service/interfaces"
)

type listeningService struct {
	numListeners     int
	numListenersLock *sync.RWMutex

	listeningChan chan interface{}
	listeningLock *sync.Mutex
}

func NewListeningService(lC chan interface{}) i.ListeningService {
	return &listeningService{
		listeningChan: lC,
	}
}

func (s *listeningService) AddListener() {
	s.numListenersLock.Lock()
	s.numListeners++
	s.numListenersLock.Unlock()
}

func (s *listeningService) RemoveListener() {
	s.numListenersLock.Lock()
	s.numListeners--
	s.numListenersLock.Unlock()
}

func (s *listeningService) Send(msg interface{}) {
	s.numListenersLock.RLock()
	s.listeningLock.Lock()
	for i := 0; i < s.numListeners; i++ {
		s.listeningChan <- msg
	}
	s.listeningLock.Unlock()
	s.numListenersLock.RUnlock()
}
