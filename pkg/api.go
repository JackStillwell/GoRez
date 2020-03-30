package gorez

import internal "github.com/JackStillwell/Gorez/internal"

// APIBase contains the information to create all requests to an API
type APIBase struct {
	baseURL        string
	devID          string
	devKey         string
	returnDataType string
	httpGet        httpGetter
}

// NewAPIBase creates an APIBase given baseURL, devID, and devKey
func NewAPIBase(baseURL, returnDataType, devID, devKey string) APIBase {
	return APIBase{
		baseURL:        baseURL,
		devID:          devID,
		devKey:         devKey,
		returnDataType: returnDataType,
		httpGet:        internal.defaultGetter{},
	}
}

func mockAPIBase(mockHTTPGetter internal.httpGetter) APIBase {
	return APIBase{
		baseURL:        "mockBaseURL",
		devID:          "mockDevID",
		devKey:         "mockDevKey",
		returnDataType: "mockReturnDataType",
		httpGet:        mockHTTPGetter,
	}
}

// Nothing returns the string "nothing"
func Nothing() string {
	return "nothing"
}
