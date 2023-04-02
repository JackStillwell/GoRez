package interfaces

import (
	"net/http"

	"github.com/JackStillwell/GoRez/internal/request/models"
	"github.com/google/uuid"
)

//go:generate mockgen --source=request.go --destination=../mocks/mock_request.go --package=mock

type HTTPGet interface {
	Get(url string) (*http.Response, error)
}

type Service interface {
	RequestManager
	Requester
}

type RequestManager interface {
	MakeRequest(*models.Request)
	GetResponse(*uuid.UUID) *models.RequestResponse
}

type Requester interface {
	Request(*models.Request) *models.RequestResponse
}
