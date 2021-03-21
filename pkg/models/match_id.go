package models

type MatchID struct {
	ActiveFlag *string `json:"Active_Flag,omitempty"`
	Match      *string `json:"Match,omitempty"`
	RetMsg     *string `json:"ret_msg,omitempty"`
}
