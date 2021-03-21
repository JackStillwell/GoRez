package interfaces

import (
	"net/http"

	"github.com/JackStillwell/GoRez/internal/request_service/models"
)

type HTTPGet interface {
	Get(url string) (*http.Response, error)
}

type RequestService interface {
	RequestManager
	Requester
}

type RequestManager interface {
	MakeRequest(*models.Request)
	GetResponse() *models.RequestResponse
	Close()
}

type Requester interface {
	Request(*models.Request) *models.RequestResponse
}
