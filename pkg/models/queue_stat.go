package models

type QueueStat struct {
	Assists    *int64  `json:"Assists,omitempty"`
	Deaths     *int64  `json:"Deaths,omitempty"`
	God        *string `json:"God,omitempty"`
	GodID      *int64  `json:"GodId,omitempty"`
	Gold       *int64  `json:"Gold,omitempty"`
	Kills      *int64  `json:"Kills,omitempty"`
	LastPlayed *string `json:"LastPlayed,omitempty"`
	Losses     *int64  `json:"Losses,omitempty"`
	Matches    *int64  `json:"Matches,omitempty"`
	Minutes    *int64  `json:"Minutes,omitempty"`
	Queue      *string `json:"Queue,omitempty"`
	Wins       *int64  `json:"Wins,omitempty"`
	PlayerID   *string `json:"player_id,omitempty"`
	RetMsg     *string `json:"ret_msg,omitempty"`
}
