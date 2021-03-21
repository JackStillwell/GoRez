package request_service

import (
	"io"
	"net/http"
	"runtime"
	"sync"

	i "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	m "github.com/JackStillwell/GoRez/internal/request_service/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type requestService struct {
	httpGet      i.HTTPGet
	requestChan  chan *m.Request
	responses    []*m.RequestResponse
	responseLock sync.Mutex
}

type httpGetter struct{}

func (_ *httpGetter) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func NewMockRequestService(capacity int, hG i.HTTPGet) i.RequestService {
	requests := make(chan *m.Request, capacity)
	responses := make([]*m.RequestResponse, capacity)

	rS := &requestService{
		httpGet:     hG,
		requestChan: requests,
		responses:   responses,
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		go requestServiceRoutine(rS)
	}

	return rS
}

func NewRequestService(capacity int) i.RequestService {
	requests := make(chan *m.Request, capacity)
	responses := make([]*m.RequestResponse, capacity)

	rS := &requestService{
		httpGet:     &httpGetter{},
		requestChan: requests,
		responses:   responses,
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		go requestServiceRoutine(rS)
	}

	return rS
}

func (s *requestService) Request(r *m.Request) *m.RequestResponse {
	requestURL, err := r.JITBuild(r.URLBuilder, r.BuildArgs)
	if err != nil {
		return &m.RequestResponse{
			Id:   r.Id,
			Resp: nil,
			Err:  errors.Wrap(err, "building requesturl"),
		}
	}

	resp, err := http.Get(requestURL)
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

// GetResponse takes a UUID and returns either the first response for a nil UUID, the response
// containing the UUID passed, or an empty response with an err if the requested response is not
// found.
func (s *requestService) GetResponse(u *uuid.UUID) (toRet *m.RequestResponse) {
	s.responseLock.Lock()
	defer s.responseLock.Unlock()

	toRet = &m.RequestResponse{
		Id:   nil,
		Resp: nil,
		Err:  errors.New("response not found"),
	}
	if u == nil {
		if len(s.responses) == 0 {
			return
		}
		toRet = s.responses[0]
		s.responses = s.responses[1:]
		return
	}

	for _, r := range s.responses {
		if r.Id == u {
			return r
		}
	}

	return
}

func requestServiceRoutine(s *requestService) {
	select {
	case r := <-s.requestChan:
		response := s.Request(r)
		s.responseLock.Lock()
		s.responses = append(s.responses, response)
		s.responseLock.Unlock()
	}
}
