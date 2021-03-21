package models

type PlayerID struct {
	PlayerID    *int64  `json:"player_id,omitempty"`
	Portal      *string `json:"portal,omitempty"`
	PortalID    *string `json:"portal_id,omitempty"`
	PrivacyFlag *string `json:"privacy_flag,omitempty"`
	RetMsg      *string `json:"ret_msg,omitempty"`
}
