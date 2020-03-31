package gorezinternal

// APIBase contains the information to create all requests to an API
type APIBase struct {
	baseURL        string
	devID          string
	devKey         string
	returnDataType string
	requester      HTTPGetter
}

// NewAPIBase creates an APIBase given baseURL, devID, and devKey
func NewAPIBase(baseURL, returnDataType, devID, devKey string) APIBase {
	return APIBase{
		baseURL:        baseURL,
		devID:          devID,
		devKey:         devKey,
		returnDataType: returnDataType,
		requester:      DefaultGetter{},
	}
}

func mockAPIBase(mockHTTPGetter HTTPGetter) APIBase {
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

// NewLimitConstants returns a LimitConstants populated with currently known values
func NewLimitConstants() LimitConstants {
	return LimitConstants{
		ConcurrentSessions: 45,
		SessionsPerDay:     500,
		SessionTimeLimit:   900, // NOTE: SessionTimeLimit is in seconds
		RequestsPerDay:     7500,
	}
}
