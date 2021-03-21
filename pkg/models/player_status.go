package models

type PlayerStatus struct {
	Match                 *int64  `json:"Match,omitempty"`
	MatchQueueID          *int64  `json:"match_queue_id,omitempty"`
	PersonalStatusMessage *string `json:"personal_status_message,omitempty"`
	RetMsg                *string `json:"ret_msg,omitempty"`
	Status                *Status `json:"status,omitempty"`
	StatusString          *string `json:"status_string,omitempty"`
}

type Status int

const (
	Offline int = iota
	InLobby
	GodSelection
	InGame
	Online
	Unknown
)

func (s Status) String() string {
	return [...]string{
		"Offline",
		"In Lobby",
		"god Selection",
		"In Game",
		"Online",
		"Unknown",
	}[s]
}
