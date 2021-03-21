package models

type Item struct {
	ActiveFlag      *string          `json:"ActiveFlag,omitempty"`
	ChildItemID     *int64           `json:"ChildItemId,omitempty"`
	DeviceName      *string          `json:"DeviceName,omitempty"`
	IconID          *int64           `json:"IconId,omitempty"`
	ItemDescription *ItemDescription `json:"ItemDescription,omitempty"`
	ItemID          *int64           `json:"ItemId,omitempty"`
	ItemTier        *int64           `json:"ItemTier,omitempty"`
	Price           *int64           `json:"Price,omitempty"`
	RestrictedRoles *string          `json:"RestrictedRoles,omitempty"`
	RootItemID      *int64           `json:"RootItemId,omitempty"`
	ShortDesc       *string          `json:"ShortDesc,omitempty"`
	StartingItem    *bool            `json:"StartingItem,omitempty"`
	Type            *string          `json:"Type,omitempty"`
	ItemIconURL     *string          `json:"itemIcon_URL,omitempty"`
	RetMsg          interface{}      `json:"ret_msg"`
}

type ItemDescription struct {
	Description          *string     `json:"Description,omitempty"`
	MenuItems            []MenuItem  `json:"Menuitems,omitempty"`
	SecondaryDescription interface{} `json:"SecondaryDescription"`
}
