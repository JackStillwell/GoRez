package auth

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/JackStillwell/GoRez/internal/base"

	i "github.com/JackStillwell/GoRez/internal/auth/interfaces"
	m "github.com/JackStillwell/GoRez/internal/auth/models"
)

type service struct {
	Auth m.Auth
	Base base.Service
}

func NewService(a m.Auth, b base.Service) i.Service {
	return &service{
		Auth: a,
		Base: b,
	}
}

func (*service) GetTimestamp(t time.Time) string {
	return t.Format("20060102150405")
}

func (s *service) GetSignature(endpoint, timestamp string) string {
	tohash := []byte(s.Auth.ID + endpoint + s.Auth.Key + timestamp)
	hash := md5.Sum(tohash)

	return hex.EncodeToString(hash[:16])
}

func (s *service) GetID() string {
	return s.Auth.ID
}
