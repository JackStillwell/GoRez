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

// MakeRequest makes a request to the HiRez API
func (t *RequestManager) MakeRequest(method string) ([]byte, error) {
	// do all the security stuff
	// construct the URL
	// use Get
	return nil, nil
}
