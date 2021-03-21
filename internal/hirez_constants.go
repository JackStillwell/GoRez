package internal

// APIConstants contains the known method fields
type APIConstants struct {
	SmiteURLBase  string
	GetDataUsed   string
	CreateSession string
	TestSession   string
}

// New returns an APIConstants obj populated with currently known values
func (t APIConstants) New() APIConstants {
	return APIConstants{
		SmiteURLBase:  "http://api.smitegame.com/smiteapi.svc",
		GetDataUsed:   "getdataused",
		CreateSession: "createsession",
		TestSession:   "testsession",
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

// New returns a SmiteConstants populated with currently known values
func (t SmiteConstants) New() SmiteConstants {
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

// New returns a ReturnDataType populated with currently known values
func (t ReturnDataType) New() ReturnDataType {
	return ReturnDataType{
		JSON: "json",
		XML:  "xml",
	}
}
