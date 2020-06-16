package gorezinternal

import "time"

// HTTPGetter is an interface for mocking
type HTTPGetter interface {
	Get(url string) ([]byte, error)
}

// RequestManagement is an interface for RequestManager
type RequestManagement interface {
	getSignature(endpoint string, timestamp string) string
	getTimestamp(currTime time.Time) string
	EndpointRequest(
		endpoint string,
		sessionID string,
		args string,
	) ([]byte, error)
	CreateSessionRequest() ([]byte, error)
}

// SessionManagement is an interface for SessionManager
type SessionManagement interface {
	mock(rm RequestManagement)
	Initialize(initFile string, rm RequestManagement)
	Save(saveFile string)
	GetSession() (string, error)
	ReturnSession(sessionID string) error
}
