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
	responseChans map[uuid.UUID]chan *m.RequestResponse
	lock          *sync.RWMutex
}

func NewService(capacity int) i.Service {
	return &service{
		responseChans: make(map[uuid.UUID]chan *m.RequestResponse, capacity),
		lock:          &sync.RWMutex{},
	}
}

func (s *service) Request(rqst *m.Request) (rr *m.RequestResponse) {
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

	if rr.Err != nil {
		rr.Err = fmt.Errorf("%s: %s : %s", rr.Err.Error(), requestURL, body)
	}

	rr.Resp = body
	return
}

func (s *service) MakeRequest(req *m.Request) {
	go func(s *service, req *m.Request) {
		log.Printf("making request %s\n", req.Id.String())
		s.ensureResponseChan(*req.Id)
		s.lock.Lock()
		s.responseChans[*req.Id] <- s.Request(req)
		s.lock.Unlock()
		log.Printf("response stored %s\n", req.Id.String())
	}(s, req)
}

func (s *service) GetResponse(id *uuid.UUID) *m.RequestResponse {
	defer func() {
		s.lock.Lock()
		delete(s.responseChans, *id)
		s.lock.Unlock()
	}()
	defer log.Printf("response returned %s\n", id.String())
	s.ensureResponseChan(*id)
	log.Printf("returning response %s\n", id.String())

	s.lock.Lock()
	retVal := <-s.responseChans[*id]
	s.lock.Unlock()
	return retVal
}

func (s *service) ensureResponseChan(id uuid.UUID) {
	s.lock.Lock()
	if _, ok := s.responseChans[id]; !ok {
		s.responseChans[id] = make(chan *m.RequestResponse, 1)
	}
	s.lock.Unlock()
}
