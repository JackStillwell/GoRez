package request_service

import (
	"io"
	"net/http"
	"runtime"

	i "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	m "github.com/JackStillwell/GoRez/internal/request_service/models"
	"github.com/pkg/errors"
)

type requestService struct {
	http         i.HTTPGet
	requestChan  chan *m.Request
	responseChan chan *m.RequestResponse
	workerKill   []chan bool
}

type httpGetter struct{}

func (*httpGetter) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func NewTestRequestService(capacity int, hG i.HTTPGet) i.RequestService {
	requests := make(chan *m.Request, capacity)
	responses := make(chan *m.RequestResponse, capacity)

	rS := &requestService{
		http:         hG,
		requestChan:  requests,
		responseChan: responses,
	}

	wKs := make([]chan bool, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		wKs[i] = make(chan bool)
		go requestServiceRoutine(rS, wKs[i])
	}

	rS.workerKill = wKs

	return rS
}

func NewRequestService(capacity int) i.RequestService {
	return NewTestRequestService(capacity, &httpGetter{})
}

func (s *requestService) Request(r *m.Request) *m.RequestResponse {
	requestURL, err := r.JITBuild(r.JITArgs)
	if err != nil {
		return &m.RequestResponse{
			Id:   r.Id,
			Resp: nil,
			Err:  errors.Wrap(err, "building requesturl"),
		}
	}

	resp, err := s.http.Get(requestURL)
	if err != nil {
		return &m.RequestResponse{
			Id:   r.Id,
			Resp: nil,
			Err:  errors.Wrap(err, "getting response"),
		}
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return &m.RequestResponse{
			Id:   r.Id,
			Resp: nil,
			Err:  errors.Wrap(err, "reading body"),
		}
	}

	return &m.RequestResponse{
		Id:   r.Id,
		Resp: body,
		Err:  nil,
	}
}

func (s *requestService) MakeRequest(r *m.Request) {
	s.requestChan <- r
}

func (s *requestService) GetResponse() (toRet *m.RequestResponse) {
	return <-s.responseChan
}

func (s *requestService) Close() {
	for _, c := range s.workerKill {
		c <- true
	}
}

func requestServiceRoutine(s *requestService, killChan chan bool) {
	kill := false
	for !kill {
		select {
		case r := <-s.requestChan:
			s.responseChan <- s.Request(r)
		case <-killChan:
			kill = true
		}
	}
}
