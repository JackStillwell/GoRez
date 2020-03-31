package gorezinternal

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

// RequestManager contains the information necessary to handle requests for HiRez API
type RequestManager struct {
	urlBase        string
	numRequests    uint16
	returnDataType string
	auth           Auth
	requester      HTTPGetter
}

// Auth contains the information necessary to authenticate against HiRez API
type Auth struct {
	devID  string
	devKey string
}

// getSignature creates the md5 signature for a request
func (t *RequestManager) getSignature(endpoint string, timestamp string) string {
	hash := md5.New()

	io.WriteString(hash, t.auth.devID)
	io.WriteString(hash, endpoint)
	io.WriteString(hash, t.auth.devKey)
	io.WriteString(hash, timestamp)

	signature := string(hash.Sum(nil))
	return signature
}

// getTimestamp creates the timestamp for a request
func getTimestamp(currTime time.Time) string {
	timestamp := currTime.Format("20060102150405")
	return timestamp
}

// EndpointRequest sends a request to the specified endpoint
func (t *RequestManager) EndpointRequest(
	endpoint string,
	sessionID string,
	args string,
) ([]byte, error) {
	// format the url properly
	timestamp := getTimestamp(time.Now().UTC())

	request := fmt.Sprintf(
		"%s/%s%s/%s/%s/%s/%s/%s",
		t.urlBase,
		endpoint,
		t.returnDataType,
		t.auth.devID,
		t.getSignature(endpoint, timestamp),
		sessionID,
		timestamp,
		args,
	)

	return t.requester.Get(request)
}

// CreateSessionRequest sends a request to the createsession endpoint
func (t *RequestManager) CreateSessionRequest() ([]byte, error) {
	// format the url properly
	timestamp := getTimestamp(time.Now().UTC())

	apiConsts := APIConstants{}.New()

	request := fmt.Sprintf(
		"%s/%s%s/%s/%s/%s",
		t.urlBase,
		apiConsts.CreateSession,
		t.returnDataType,
		t.auth.devID,
		t.getSignature(apiConsts.CreateSession, timestamp),
		timestamp,
	)

	return t.requester.Get(request)
}
