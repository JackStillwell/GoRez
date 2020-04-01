package gorezinternal

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
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

// mock creates a mock RequestManager for testing purposes
func (t RequestManager) mock(mockRequester HTTPGetter) RequestManager {
	return RequestManager{
		urlBase:        "mockURLBase",
		numRequests:    0,
		returnDataType: "json",
		auth: Auth{
			devID:  "mockDevID",
			devKey: "mockDevKey",
		},
		requester: mockRequester,
	}
}

// Auth contains the information necessary to authenticate against HiRez API
type Auth struct {
	devID  string
	devKey string
}

// GetSignature creates the md5 signature for a request
func (t *RequestManager) GetSignature(endpoint string, timestamp string) string {
	tohash := []byte(t.auth.devID + endpoint + t.auth.devKey + timestamp)
	hash := md5.Sum(tohash)

	return hex.EncodeToString(hash[:16])
}

// GetTimestamp creates the timestamp for a request
func GetTimestamp(currTime time.Time) string {
	timestamp := currTime.Format("20060102150405")
	return timestamp
}

// EndpointRequest sends a request to the specified endpoint
func (t *RequestManager) EndpointRequest(
	endpoint string,
	sessionID string,
	args string,
	timestampTime time.Time,
) ([]byte, error) {
	timestamp := GetTimestamp(timestampTime)

	// format the url properly
	request := fmt.Sprintf(
		"%s/%s%s/%s/%s/%s/%s",
		t.urlBase,
		endpoint,
		t.returnDataType,
		t.auth.devID,
		t.GetSignature(endpoint, timestamp),
		sessionID,
		timestamp,
	)

	if args != "" {
		request += "/" + args
	}

	return t.requester.Get(request)
}

// CreateSessionRequest sends a request to the createsession endpoint
func (t *RequestManager) CreateSessionRequest() ([]byte, error) {
	// format the url properly
	timestamp := GetTimestamp(time.Now().UTC())

	apiConsts := APIConstants{}.New()

	request := fmt.Sprintf(
		"%s/%s%s/%s/%s/%s",
		t.urlBase,
		apiConsts.CreateSession,
		t.returnDataType,
		t.auth.devID,
		t.GetSignature(apiConsts.CreateSession, timestamp),
		timestamp,
	)

	return t.requester.Get(request)
}
