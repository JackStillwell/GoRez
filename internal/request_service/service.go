package request_service

import (
	"io"
	"net/http"

	i "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	m "github.com/JackStillwell/GoRez/internal/request_service/models"
	"github.com/pkg/errors"
)

type requestService struct {
	Requests     []m.Request
	Responses    []m.RequestResponse
	requestChan  chan m.Request
	responseChan chan m.RequestResponse
}

func (s *requestService) NewRequestService(capacity int) i.RequestService {
	requests := make([]m.Request, capacity)
	responses := make([]m.RequestResponse, capacity)
	return &requestService{
		Requests:  requests,
		Responses: responses,
	}
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
