package models

type UsageInfo struct {
	ActiveSessions     *int64  `json:"Active_Sessions,omitempty"`
	ConcurrentSessions *int64  `json:"Concurrent_Sessions,omitempty"`
	RequestLimitDaily  *int64  `json:"Request_Limit_Daily,omitempty"`
	SessionCap         *int64  `json:"Session_Cap,omitempty"`
	SessionTimeLimit   *int64  `json:"Session_Time_Limit,omitempty"`
	TotalRequestsToday *int64  `json:"Total_Requests_Today,omitempty"`
	TotalSessionsToday *int64  `json:"Total_Sessions_Today,omitempty"`
	RetMsg             *string `json:"ret_msg,omitempty"`
}
