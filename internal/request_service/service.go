package request_service

import (
	"io"
	"net/http"
	"sync"

	i "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	m "github.com/JackStillwell/GoRez/internal/request_service/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type requestService struct {
	i.Requester
	i.RequestManager
}

type requester struct {
	http i.HTTPGet
}

type requestManager struct {
	r i.Requester

	responses     []*m.RequestResponse
	responsesLock *sync.RWMutex

	freeNotify chan int

	listenerNotify chan int
	listenerCount  int
	listenerLock   *sync.Mutex
}

type httpGetter struct{}

func (*httpGetter) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func NewTestRequester(http i.HTTPGet) i.Requester {
	return &requester{http}
}

func NewTestRequestManager(capacity int, r i.Requester) *requestManager {
	responses := make([]*m.RequestResponse, capacity)

	freeNotifyChan := make(chan int, capacity)
	for i := 0; i < capacity; i++ {
		freeNotifyChan <- i
	}

	listenerChan := make(chan int)

	rM := &requestManager{
		r: r,

		responses:     responses,
		responsesLock: &sync.RWMutex{},

		freeNotify: freeNotifyChan,

		listenerNotify: listenerChan,
		listenerCount:  0,
		listenerLock:   &sync.Mutex{},
	}

	return rM
}

func NewRequestService(capacity int) i.RequestService {
	return &requestService{
		Requester:      NewTestRequester(&httpGetter{}),
		RequestManager: NewTestRequestManager(capacity, NewTestRequester(&httpGetter{})),
	}
}

func (r *requester) Request(rqst *m.Request) (rr *m.RequestResponse) {
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

	resp, err := r.http.Get(requestURL)
	if err != nil {
		rr.Err = errors.Wrap(err, "getting response")
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		rr.Err = errors.Wrap(err, "reading body")
		return
	}

	rr.Resp = body
	return
}

func (rM *requestManager) MakeRequest(r *m.Request) {
	go rM.request(r)
}

func (rM *requestManager) GetResponse(id *uuid.UUID, retChan chan *m.RequestResponse) {
	rM.listenerLock.Lock()
	result := rM.searchResponse(id)

	if result == nil {
		rM.listenerCount++
	}

	rM.listenerLock.Unlock()

	for result == nil {
		idx := <-rM.listenerNotify
		if rM.responses[idx].Id == id {
			rM.decrementListeners()
			defer freeIdx(rM, idx)
			result = rM.responses[idx]
		}
	}

	retChan <- result
}

func (rM *requestManager) searchResponse(id *uuid.UUID) *m.RequestResponse {
	rM.responsesLock.RLock()
	for idx, v := range rM.responses {
		if v != nil && v.Id == id {
			rM.responsesLock.RUnlock()
			defer freeIdx(rM, idx)
			return rM.responses[idx]
		}
	}
	rM.responsesLock.RUnlock()

	return nil
}

func freeIdx(rM *requestManager, idx int) {
	rM.updateResponses(idx, nil)
	rM.freeNotify <- idx
}

func (rM *requestManager) request(rqst *m.Request) {
	response := rM.r.Request(rqst)
	responseIdx := <-rM.freeNotify
	rM.updateResponses(responseIdx, response)
	rM.notifyListeners(responseIdx)
}

func (rM *requestManager) updateResponses(idx int, v *m.RequestResponse) {
	rM.responsesLock.Lock()
	rM.responses[idx] = v
	rM.responsesLock.Unlock()
}

func (rM *requestManager) notifyListeners(idx int) {
	rM.listenerLock.Lock()
	for i := 0; i < rM.listenerCount; i++ {
		rM.listenerNotify <- idx
	}
	rM.listenerLock.Unlock()
}

func (rM *requestManager) decrementListeners() {
	rM.listenerLock.Lock()
	rM.listenerCount--
	rM.listenerLock.Unlock()
}
