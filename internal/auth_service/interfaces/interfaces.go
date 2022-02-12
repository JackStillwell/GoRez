package interfaces

import "time"

//go:generate mockgen --source=interfaces.go --destination=../mocks/mock_auth_service.go --package=mocks
type AuthService interface {
	GetTimestamp(t time.Time) string
	GetSignature(endpoint, timestamp string) string
	GetID() string
}
