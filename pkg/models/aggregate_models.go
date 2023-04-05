package models

type PlayerIDWithName struct {
	PlayerID `json:",inline"`
	Name     string `json:"Name"`
}

type MatchIDWithQueue struct {
	MatchID `json:",inline"`
	QueueID int `json:"QueueID"`
}
