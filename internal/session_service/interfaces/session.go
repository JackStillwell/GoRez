package interfaces

import m "github.com/JackStillwell/GoRez/internal/session_service/models"

type SessionService interface {
	ReserveSession(int, chan *m.Session)
	ReleaseSession([]*m.Session)
	BadSession([]*m.Session)
}
