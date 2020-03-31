package gorezinternal

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
