package gorez

import internal "github.com/JackStillwell/GoRez/internal"

// APIBase contains the information to create all requests to an API
type APIBase struct {
	baseURL        string
	devID          string
	devKey         string
	returnDataType string
	httpGet        internal.HTTPGetter
}

// NewAPIBase creates an APIBase given baseURL, devID, and devKey
func NewAPIBase(baseURL, returnDataType, devID, devKey string) APIBase {
	return APIBase{
		baseURL:        baseURL,
		devID:          devID,
		devKey:         devKey,
		returnDataType: returnDataType,
		httpGet:        internal.DefaultGetter{},
	}
}

func mockAPIBase(mockHTTPGetter internal.HTTPGetter) APIBase {
	return APIBase{
		baseURL:        "mockBaseURL",
		devID:          "mockDevID",
		devKey:         "mockDevKey",
		returnDataType: "mockReturnDataType",
		httpGet:        mockHTTPGetter,
	}
}

// GetSession returns a session id
func GetSession(api APIBase) (string, error) {
	return internal.GetSession(
		api.baseURL,
		api.returnDataType,
		api.devID,
		api.devKey,
		api.httpGet,
	)
}
