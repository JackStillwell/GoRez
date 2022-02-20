package interfaces

import (
	sessionM "github.com/JackStillwell/GoRez/internal/session_service/models"

	requestM "github.com/JackStillwell/GoRez/internal/request_service/models"
)

type GorezUtil interface {
	BulkAsyncSessionRequest([]func(*sessionM.Session) *requestM.Request) ([][]byte, []error)
	MultiRequest(requestArgs []string, endpoint, method string) ([][]byte, []error)
	SingleRequest(r requestM.Request, unmarshalTo interface{}) error
}
