package interfaces

import (
	sessionM "github.com/JackStillwell/GoRez/internal/session_service/models"

	requestM "github.com/JackStillwell/GoRez/internal/request_service/models"
)

//go:generate mockgen --source=./gorez_util.go --destination=../mocks/mock_gorez_util.go --package=mock
type GorezUtil interface {
	BulkAsyncSessionRequest([]func(*sessionM.Session) *requestM.Request) ([][]byte, []error)
	MultiRequest(requestArgs []string, endpoint, method string) ([][]byte, []error)
	SingleRequest(url, endpoint, endpointArgs string) ([]byte, error)
}
