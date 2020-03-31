package gorezinternal

// RequestManager contains the information necessary to handle requests for HiRez API
type RequestManager struct {
	urlBase        string
	numRequests    uint16
	returnDataType string
	auth           Auth
}

// Auth contains the information necessary to authenticate against HiRez API
type Auth struct {
	devID  string
	devKey string
}

// makeRequest makes a request to the HiRez API
func (t *RequestManager) makeRequest(endpoint string) ([]byte, error) {
	// do all the security stuff
	// construct the URL
	// use Get
	return nil, nil
}

// EndpointRequest sends a request to the specified endpoint
func (t *RequestManager) EndpointRequest(endpoint string) ([]byte, error) {
	// format the url properly

	return t.makeRequest(endpoint)
}

// CreateSessionRequest sends a request to the createsession endpoint
func (t *RequestManager) CreateSessionRequest() ([]byte, error) {
	// format the url properly

	return t.makeRequest("createsession")
}
