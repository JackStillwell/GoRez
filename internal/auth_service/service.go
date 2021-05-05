package auth_service

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	i "github.com/JackStillwell/GoRez/internal/auth_service/interfaces"
	m "github.com/JackStillwell/GoRez/internal/auth_service/models"
)

type authService struct {
	Auth m.Auth
}

func NewAuthService(a m.Auth) i.AuthService {
	return &authService{
		Auth: a,
	}
}

func (r *authService) GetTimestamp(t time.Time) string {
	timestamp := t.Format("20060102150405")
	return timestamp
}

func (r *authService) GetSignature(endpoint, timestamp string) string {
	tohash := []byte(r.Auth.ID + endpoint + r.Auth.Key + timestamp)
	hash := md5.Sum(tohash)

	return hex.EncodeToString(hash[:16])
}

func (r *authService) GetID() string {
	return r.Auth.ID
}
