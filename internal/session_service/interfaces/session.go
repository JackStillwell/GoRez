package interfaces

import m "github.com/JackStillwell/GoRez/internal/session_service/models"

type SessionService interface {
	ReserveSession(int) ([]*m.Session, error)
	ReleaseSession([]*m.Session)
	BadSession([]*m.Session)
}
