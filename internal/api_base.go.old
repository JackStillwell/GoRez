package internal

// APIBase contains the information to create all requests to an API
type APIBase struct {
	baseURL        string
	devID          string
	devKey         string
	returnDataType string
	requester      HTTPGetter
}

// New creates an APIBase given baseURL, devID, and devKey
func (t APIBase) New(baseURL, returnDataType, devID, devKey string) APIBase {
	return APIBase{
		baseURL:        baseURL,
		devID:          devID,
		devKey:         devKey,
		returnDataType: returnDataType,
		requester:      DefaultGetter{},
	}
}

// mock creates a mock APIBase given a mock getter
func (t APIBase) mock(mockHTTPGetter HTTPGetter) APIBase {
	return APIBase{
		baseURL:        "mockBaseURL",
		devID:          "mockDevID",
		devKey:         "mockDevKey",
		returnDataType: "mockReturnDataType",
		requester:      mockHTTPGetter,
	}
}

// LimitConstants contains the known limit constants
type LimitConstants struct {
	ConcurrentSessions uint8
	SessionsPerDay     uint16
	SessionTimeLimit   uint16
	RequestsPerDay     uint16
}

// New returns a LimitConstants populated with currently known values
func (t LimitConstants) New() LimitConstants {
	return LimitConstants{
		ConcurrentSessions: 45,
		SessionsPerDay:     500,
		SessionTimeLimit:   900, // NOTE: SessionTimeLimit is in seconds
		RequestsPerDay:     7500,
	}
}
