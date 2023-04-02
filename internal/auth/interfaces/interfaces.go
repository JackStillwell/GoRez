package interfaces

import "time"

//go:generate mockgen --source=interfaces.go --destination=../mocks/mock_service.go --package=mocks
type Service interface {
	GetTimestamp(t time.Time) string
	GetSignature(endpoint, timestamp string) string
	GetID() string
}
