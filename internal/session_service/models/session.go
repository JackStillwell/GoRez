package models

import "time"

type Session struct {
	Key     string
	Created *time.Time
}
