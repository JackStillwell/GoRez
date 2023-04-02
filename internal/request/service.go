package request

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	i "github.com/JackStillwell/GoRez/internal/request/interfaces"
	m "github.com/JackStillwell/GoRez/internal/request/models"
)

type service struct {
	responses       []*m.RequestResponse
	responseIdToIdx map[*uuid.UUID]int
	responsesLock   *sync.RWMutex

	freeNotify chan int

	listenerNotify chan int

	listenerCount    int
	numListenersLock *sync.RWMutex
}

func NewService(capacity int) i.Service {
	responses := make([]*m.RequestResponse, capacity)

	freeNotifyChan := make(chan int, capacity)
	for i := 0; i < capacity; i++ {
		freeNotifyChan <- i
	}

	listenerChan := make(chan int)

	return &service{
		responses:       responses,
		responseIdToIdx: make(map[*uuid.UUID]int, capacity),
		responsesLock:   &sync.RWMutex{},

		freeNotify: freeNotifyChan,

		listenerNotify: listenerChan,

		listenerCount:    0,
		numListenersLock: &sync.RWMutex{},
	}
}

func (s *service) Request(rqst *m.Request) (rr *m.RequestResponse) {
	log.Println("making request")

	rr = &m.RequestResponse{
		Id:   rqst.Id,
		Resp: nil,
		Err:  nil,
	}

	requestURL, err := rqst.JITFunc()
	if err != nil {
		rr.Err = errors.Wrap(err, "building requesturl")
		return
	}
	log.Println("request url:", requestURL)

	resp, err := http.Get(requestURL)
	if err != nil {
		rr.Err = errors.Wrap(err, "getting response")
		return
	}

	if resp.StatusCode != http.StatusOK {
		rr.Err = fmt.Errorf("status code %d, require %d", resp.StatusCode, http.StatusOK)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		rr.Err = errors.Wrap(err, "reading body")
		return
	}
	log.Println("response received:", requestURL)

	if rr.Err != nil {
		rr.Err = fmt.Errorf("%s: %s", rr.Err.Error(), body)
	}

	rr.Resp = body
	return
}

func (s *service) MakeRequest(req *m.Request) {
	go func(r *service, req *m.Request) {
		// defer ginkgo.GinkgoRecover()
		r.request(req)
	}(s, req)
}

func (s *service) GetResponse(id *uuid.UUID, retChan chan *m.RequestResponse) {
	log.Println("locking request responses")
	s.responsesLock.RLock()
	log.Println("request responses locked")
	log.Println("searching request responses for id:", id.String())
	for idx, v := range s.responses {
		if v != nil && v.Id == id {
			log.Println("request response for id:", id.String(), "found")
			s.responseIdToIdx[id] = idx
			retChan <- s.responses[idx]
			log.Println("request response for id:", id.String(), "returned")
			s.responsesLock.RUnlock()
			log.Println("request responses unlocked")
			return
		}
	}
	log.Println("request response for id:", id.String(), "not found")

	log.Println("adding listener for new responses")
	s.numListenersLock.Lock()
	s.listenerCount++
	s.numListenersLock.Unlock()
	log.Println("listener for new responses added")
	s.responsesLock.RUnlock()
	log.Println("request responses unlocked")

	log.Println("listening for new responses")
	for {
		idx := <-s.listenerNotify
		log.Println("notified of new response")
		if s.responses[idx].Id == id {
			log.Println("request response for id:", id.String(), "found")
			log.Println("removing listener for new responses")
			s.numListenersLock.Lock()
			s.listenerCount--
			s.numListenersLock.Unlock()
			log.Println("listener for new responses removed ")
			s.responseIdToIdx[id] = idx
			retChan <- s.responses[idx]
			log.Println("request response for id:", id.String(), "returned")
			return
		}
	}
}

func (s *service) FreeResponse(id *uuid.UUID) {
	if id == nil {
		log.Println("cannot free 'nil' id response")
	}
	idx := s.responseIdToIdx[id]
	log.Printf("freeing response idx %d\n", idx)
	s.updateResponses(idx, nil)
	s.freeNotify <- idx
}

func (s *service) request(rqst *m.Request) {
	response := s.Request(rqst)

	log.Println("waiting for free slot in response buffer")
	responseIdx := <-s.freeNotify

	log.Printf("updating free slot %d in response buffer\n", responseIdx)
	s.updateResponses(responseIdx, response)
	log.Printf("slot %d in response buffer updated\n", responseIdx)
	s.notifyListeners(responseIdx)
	log.Printf("listeners notified for response buffer loc %d \n", responseIdx)
}

func (s *service) updateResponses(idx int, v *m.RequestResponse) {
	s.responsesLock.Lock()
	s.responses[idx] = v
	s.responsesLock.Unlock()
}

func (s *service) notifyListeners(idx int) {
	s.numListenersLock.RLock()
	for i := 0; i < s.listenerCount; i++ {
		s.listenerNotify <- idx
	}
	s.numListenersLock.RUnlock()
}
