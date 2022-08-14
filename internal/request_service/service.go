package request_service

import (
	"fmt"
	"io"
	"log"
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
	log.Println("making request")

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

func (r *requestService) MakeRequest(req *m.Request) {
	go func(r *requestService, req *m.Request) {
		// defer ginkgo.GinkgoRecover()
		r.request(req)
	}(r, req)
}

func (r *requestService) GetResponse(id *uuid.UUID, retChan chan *m.RequestResponse) {
	log.Println("locking request responses")
	r.responsesLock.RLock()
	log.Println("request responses locked")
	log.Println("searching request responses for id:", id.String())
	for idx, v := range r.responses {
		if v != nil && v.Id == id {
			log.Println("request response for id:", id.String(), "found")
			defer freeIdx(r, idx)
			retChan <- r.responses[idx]
			log.Println("request response for id:", id.String(), "returned")
			r.responsesLock.RUnlock()
			log.Println("request responses unlocked")
			return
		}
	}
	log.Println("request response for id:", id.String(), "not found")

	log.Println("adding listener for new responses")
	r.numListenersLock.Lock()
	r.listenerCount++
	r.numListenersLock.Unlock()
	log.Println("listener for new responses added")
	r.responsesLock.RUnlock()
	log.Println("request responses unlocked")

	log.Println("listening for new responses")
	for {
		idx := <-r.listenerNotify
		log.Println("notified of new response")
		if r.responses[idx].Id == id {
			log.Println("request response for id:", id.String(), "found")
			log.Println("removing listener for new responses")
			r.numListenersLock.Lock()
			r.listenerCount--
			r.numListenersLock.Unlock()
			log.Println("listener for new responses removed ")
			defer freeIdx(r, idx)
			retChan <- r.responses[idx]
			log.Println("request response for id:", id.String(), "returned")
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

	log.Println("waiting for free slot in response buffer")
	responseIdx := <-r.freeNotify

	log.Printf("updating free slot %d in response buffer\n", responseIdx)
	r.updateResponses(responseIdx, response)
	log.Printf("slot %d in response buffer updated\n", responseIdx)
	r.notifyListeners(responseIdx)
	log.Printf("listeners notified for response buffer loc %d \n", responseIdx)
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
