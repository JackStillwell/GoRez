package request

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/JackStillwell/GoRez/internal/base"
	i "github.com/JackStillwell/GoRez/internal/request/interfaces"
	m "github.com/JackStillwell/GoRez/internal/request/models"
)

type service struct {
	responseChans map[uuid.UUID]chan *m.RequestResponse
	lock          *sync.Mutex
	base          base.Service
}

func (s *service) delResponseChan(id uuid.UUID) {
	log := s.base.GetLogger()
	log.Debug("waiting for response chan lock", zap.String("id", id.String()))
	s.lock.Lock()
	log.Debug("response chan lock acquired", zap.String("id", id.String()))
	delete(s.responseChans, id)
	log.Debug("releasing response chan lock", zap.String("id", id.String()))
	s.lock.Unlock()
}

func (s *service) getResponseChan(id uuid.UUID) chan *m.RequestResponse {
	log := s.base.GetLogger()
	log.Debug("waiting for response chan lock", zap.String("id", id.String()))
	s.lock.Lock()
	log.Debug("response chan lock acquired", zap.String("id", id.String()))
	retVal, ok := s.responseChans[id]
	if !ok {
		retVal = make(chan *m.RequestResponse, 1)
		s.responseChans[id] = retVal
	}
	log.Debug("releasing response chan lock", zap.String("id", id.String()))
	s.lock.Unlock()
	return retVal
}

func NewService(capacity int, b base.Service) i.Service {
	return &service{
		responseChans: make(map[uuid.UUID]chan *m.RequestResponse, capacity),
		lock:          &sync.Mutex{},
		base:          b,
	}
}

func (s *service) Request(rqst *m.Request) (rr *m.RequestResponse) {
	log := s.base.GetLogger()

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

	log.Info("making request", zap.String("id", rqst.Id.String()), zap.String("url", requestURL))
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
	log.Debug("request complete", zap.String("id", rqst.Id.String()),
		zap.String("body", string(body)))
	return
}

func (s *service) MakeRequest(req *m.Request) {
	log := s.base.GetLogger()
	go func(s *service, req *m.Request) {
		log.Debug("queueing request", zap.String("id", req.Id.String()))
		respChan := s.getResponseChan(*req.Id)
		respChan <- s.Request(req)
		log.Debug("response stored", zap.String("id", req.Id.String()))
	}(s, req)
}

func (s *service) GetResponse(id *uuid.UUID) *m.RequestResponse {
	log := s.base.GetLogger()
	defer func() {
		s.delResponseChan(*id)
	}()
	defer log.Debug("response returned", zap.String("id", id.String()))

	respChan := s.getResponseChan(*id)
	return <-respChan
}
