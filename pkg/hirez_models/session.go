package hirez_models

type Session struct {
	RetMsg    *string `json:"ret_msg,omitempty"`
	SessionID *string `json:"session_id,omitempty"`
	Timestamp *string `json:"timestamp,omitempty"`
}
