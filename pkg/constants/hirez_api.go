package constants

const (
	// API Base URLs
	SmiteURLBase = "https://api.smitegame.com/smiteapi.svc"

	// API Endpoints
	GetDataUsed            = "getdataused"
	CreateSession          = "createsession"
	TestSession            = "testsession"
	GetMatchDetails        = "getmatchdetails"
	GetMatchDetailsBatch   = "getmatchdetailsbatch"
	GetMatchIDsByQueue     = "getmatchidsbyqueue"
	GetGods                = "getgods"
	GetItems               = "getitems"
	GetGodRecommendedItems = "getgodrecommendeditems"
)

type HiRezConstants struct {
	SmiteURLBase, GetDataUsed, CreateSession, TestSession, GetMatchDetails, GetMatchDetailsBatch,
	GetMatchIDsByQueue, GetGods, GetItems, GetGodRecommendedItems string
}

func NewHiRezConstants() HiRezConstants {
	return HiRezConstants{
		SmiteURLBase,
		GetDataUsed,
		CreateSession,
		TestSession,
		GetMatchDetails,
		GetMatchDetailsBatch,
		GetMatchIDsByQueue,
		GetGods,
		GetItems,
		GetGodRecommendedItems,
	}
}
