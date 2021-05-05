package interfaces

import "time"

type AuthService interface {
	GetTimestamp(t time.Time) string
	GetSignature(endpoint, timestamp string) string
	GetID() string
}
