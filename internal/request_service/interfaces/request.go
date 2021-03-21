package interfaces

import (
	"net/http"

	"github.com/JackStillwell/GoRez/internal/request_service/models"
	"github.com/google/uuid"
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
	GetResponse(*uuid.UUID) *models.RequestResponse
	Close()
}

type Requester interface {
	Request(*models.Request) *models.RequestResponse
}
