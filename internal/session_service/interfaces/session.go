package interfaces

import m "github.com/JackStillwell/GoRez/internal/session_service/models"

//go:generate mockgen --source=session.go --destination=../mocks/mock_session.go --package=mock
type SessionService interface {
	ReserveSession(int, chan *m.Session)
	ReleaseSession([]*m.Session)
	BadSession([]*m.Session)
}
