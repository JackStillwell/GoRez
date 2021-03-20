package interfaces

import (
	"github.com/JackStillwell/GoRez/internal/request_service/models"
	"github.com/google/uuid"
)

type RequestService interface {
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
