package models

type ItemRecommendation struct {
	Category        *string `json:"Category,omitempty"`
	Item            *string `json:"Item,omitempty"`
	Role            *string `json:"Role,omitempty"`
	CategoryValueID *int64  `json:"category_value_id,omitempty"`
	GodID           *int64  `json:"god_id,omitempty"`
	GodName         *string `json:"god_name,omitempty"`
	IconID          *int64  `json:"icon_id,omitempty"`
	ItemID          *int64  `json:"item_id,omitempty"`
	RetMsg          *string `json:"ret_msg,omitempty"`
	RoleValueID     *int64  `json:"role_value_id,omitempty"`
}
