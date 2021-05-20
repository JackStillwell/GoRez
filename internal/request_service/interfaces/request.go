package interfaces

import (
	"net/http"
	"time"

	"github.com/JackStillwell/GoRez/internal/request_service/models"
	"github.com/google/uuid"
)

// mockgen --source=interfaces/request.go --destination=mocks/mock_request.go --package=mock

type HTTPGet interface {
	Get(url string) (*http.Response, error)
}

type RequestService interface {
	RequestManager
	Requester
}

type RequestManager interface {
	MakeRequest(*models.Request)
	GetResponse(*uuid.UUID, chan *models.RequestResponse, time.Duration) error
	Close()
}

type Requester interface {
	Request(*models.Request) *models.RequestResponse
}
