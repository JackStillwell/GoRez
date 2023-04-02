package interfaces

import m "github.com/JackStillwell/GoRez/internal/session/models"

//go:generate mockgen --source=session.go --destination=../mocks/mock_session.go --package=mock
type Service interface {
	GetAvailableSessions() []*m.Session
	ReserveSession(int, chan *m.Session)
	ReleaseSession([]*m.Session)
	BadSession([]*m.Session)
}
