package gorez

// APIConstants contains the known method fields
type APIConstants struct {
	SmiteURLBase  string
	GetDataUsed   string
	CreateSession string
}

// NewAPIConstants returns a SmiteURL populated with currently known values
func NewAPIConstants() APIConstants {
	return APIConstants{
		SmiteURLBase:  "http://api.smitegame.com/smiteapi.svc",
		GetDataUsed:   "getdataused",
		CreateSession: "createsession",
	}
}

// SmiteConstants contains the known method fields
type SmiteConstants struct {
	GetMatchDetails      string
	GetMatchDetailsBatch string
	GetMatchIDsByQueue   string
	GetGods              string
	GetItems             string
	RankedConquest       string
	RankedJoust          string
	RankedDuel           string
}

// NewSmiteConstants returns a SmiteURL populated with currently known values
func NewSmiteConstants() SmiteConstants {
	return SmiteConstants{
		GetMatchDetails:      "getmatchdetails",
		GetMatchDetailsBatch: "getmatchdetailsbatch",
		GetMatchIDsByQueue:   "getmatchidsbyqueue",
		GetGods:              "getgods",
		GetItems:             "getitems",
		RankedConquest:       "451",
		RankedJoust:          "450",
		RankedDuel:           "440",
	}
}

// ReturnDataType contains the known return data types
type ReturnDataType struct {
	JSON string
	XML  string
}

// NewReturnDataType returns a ReturnDataType populated with currently known values
func NewReturnDataType() ReturnDataType {
	return ReturnDataType{
		JSON: "json",
		XML:  "xml",
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
