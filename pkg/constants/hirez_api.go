package constants

const (
	// API Base URLs
	SmiteURLBase = "https://api.smitegame.com/smiteapi.svc"

	// API Endpoints
	GetDataUsed            = "getdataused"
	CreateSession          = "createsession"
	TestSession            = "testsession"
	GetMatchHistory        = "getmatchhistory"
	GetMatchDetails        = "getmatchdetails"
	GetMatchDetailsBatch   = "getmatchdetailsbatch"
	GetMatchIDsByQueue     = "getmatchidsbyqueue"
	GetGods                = "getgods"
	GetItems               = "getitems"
	GetGodRecommendedItems = "getgodrecommendeditems"
	GetPlayerIDByName      = "getplayeridbyname"
)

type HiRezConstants struct {
	SmiteURLBase, GetDataUsed, CreateSession, TestSession, GetMatchHistory, GetMatchDetails, GetMatchDetailsBatch,
	GetMatchIDsByQueue, GetGods, GetItems, GetGodRecommendedItems, GetPlayerIDByName string
}

func NewHiRezConstants() HiRezConstants {
	return HiRezConstants{
		SmiteURLBase,
		GetDataUsed,
		CreateSession,
		TestSession,
		GetMatchHistory,
		GetMatchDetails,
		GetMatchDetailsBatch,
		GetMatchIDsByQueue,
		GetGods,
		GetItems,
		GetGodRecommendedItems,
		GetPlayerIDByName,
	}
}
