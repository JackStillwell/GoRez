package request_service

import (
	"io"
	"net/http"
	"runtime"
	"time"

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
	r           i.Requester
	requestChan chan *m.Request
	responses   []*m.RequestResponse
	freeNotify  chan int
	workerKill  []chan bool
}

type httpGetter struct{}

func (*httpGetter) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func NewTestRequester(http i.HTTPGet) i.Requester {
	return &requester{http}
}

func NewTestRequestService(capacity int, r i.Requester) i.RequestService {
	requests := make(chan *m.Request, capacity)
	responses := make([]*m.RequestResponse, capacity)

	freeNotifyChan := make(chan int, capacity)
	for i := 0; i < capacity; i++ {
		freeNotifyChan <- i
	}

	rM := &requestManager{
		r:           r,
		requestChan: requests,
		responses:   responses,
		freeNotify:  freeNotifyChan,
	}

	wKs := make([]chan bool, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		wKs[i] = make(chan bool)
		go requestServiceRoutine(rM, wKs[i])
	}

	rM.workerKill = wKs

	rS := &requestService{r, rM}

	return rS
}

func NewRequestService(capacity int) i.RequestService {
	return NewTestRequestService(capacity, NewTestRequester(&httpGetter{}))
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
	rM.requestChan <- r
}

func (rM *requestManager) GetResponse(id *uuid.UUID, retChan chan *m.RequestResponse, timeout time.Duration) error {
	// NOTE: this could be cleaner and less processor intensive. Rather than just looping for the timeout,
	// instead make a channel where all newly filled uuids can be checked and returned without being stored, with only
	// one loop in the beginning.
	startTime := time.Now()
	for time.Now().Sub(startTime) < timeout {
		for idx, v := range rM.responses {
			if v != nil && v.Id == id {
				defer freeIdx(rM, idx)
				retChan <- v
				return nil
			}
		}
	}
	return errors.New("timeout")
}

func (rM *requestManager) Close() {
	for _, c := range rM.workerKill {
		c <- true
	}
}

func freeIdx(rM *requestManager, idx int) {
	rM.freeNotify <- idx
}

func requestServiceRoutine(rM *requestManager, killChan chan bool) {
	kill := false
	for !kill {
		select {
		case rqst := <-rM.requestChan:
			response := rM.r.Request(rqst)
			responseIdx := <-rM.freeNotify
			rM.responses[responseIdx] = response
		case <-killChan:
			kill = true
		}
	}
}
