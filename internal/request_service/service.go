package request_service

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	i "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	m "github.com/JackStillwell/GoRez/internal/request_service/models"
)

type requestService struct {
	responses     []*m.RequestResponse
	responsesLock *sync.RWMutex

	freeNotify chan int

	listenerNotify chan int

	listenerCount    int
	numListenersLock *sync.RWMutex
}

func NewRequestService(capacity int) i.RequestService {
	responses := make([]*m.RequestResponse, capacity)

	freeNotifyChan := make(chan int, capacity)
	for i := 0; i < capacity; i++ {
		freeNotifyChan <- i
	}

	listenerChan := make(chan int)

	return &requestService{
		responses:     responses,
		responsesLock: &sync.RWMutex{},

		freeNotify: freeNotifyChan,

		listenerNotify: listenerChan,

		listenerCount:    0,
		numListenersLock: &sync.RWMutex{},
	}
}

func (r *requestService) Request(rqst *m.Request) (rr *m.RequestResponse) {
	rr = &m.RequestResponse{
		Id:   rqst.Id,
		Resp: nil,
		Err:  nil,
	}

	requestURL, err := rqst.JITBuild(rqst.JITArgs)
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
		rr.Err = fmt.Errorf("%s: %s", rr.Err.Error(), body)
	}

	rr.Resp = body
	return
}

func (r *requestService) MakeRequest(req *m.Request) {
	go func(r *requestService, req *m.Request) {
		// defer ginkgo.GinkgoRecover()
		r.request(req)
	}(r, req)
}

func (r *requestService) GetResponse(id *uuid.UUID, retChan chan *m.RequestResponse) {
	r.responsesLock.RLock()
	for idx, v := range r.responses {
		if v != nil && v.Id == id {
			r.responsesLock.RUnlock()
			defer freeIdx(r, idx)
			retChan <- r.responses[idx]
			return
		}
	}

	r.numListenersLock.Lock()
	r.listenerCount++
	r.numListenersLock.Unlock()
	r.responsesLock.RUnlock()

	for {
		idx := <-r.listenerNotify
		if r.responses[idx].Id == id {
			r.numListenersLock.Lock()
			r.listenerCount--
			r.numListenersLock.Unlock()
			defer freeIdx(r, idx)
			retChan <- r.responses[idx]
			return
		}
	}
}

func freeIdx(r *requestService, idx int) {
	r.updateResponses(idx, nil)
	r.freeNotify <- idx
}

func (r *requestService) request(rqst *m.Request) {
	response := r.Request(rqst)
	responseIdx := <-r.freeNotify
	r.updateResponses(responseIdx, response)
	r.notifyListeners(responseIdx)
}

func (r *requestService) updateResponses(idx int, v *m.RequestResponse) {
	r.responsesLock.Lock()
	r.responses[idx] = v
	r.responsesLock.Unlock()
}

func (r *requestService) notifyListeners(idx int) {
	r.numListenersLock.RLock()
	for i := 0; i < r.listenerCount; i++ {
		r.listenerNotify <- idx
	}
	r.numListenersLock.RUnlock()
}
