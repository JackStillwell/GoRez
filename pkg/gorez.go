package gorez

import (
	c "github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"

	auth "github.com/JackStillwell/GoRez/internal/auth_service"
	authI "github.com/JackStillwell/GoRez/internal/auth_service/interfaces"
	authM "github.com/JackStillwell/GoRez/internal/auth_service/models"

	request "github.com/JackStillwell/GoRez/internal/request_service"
	requestI "github.com/JackStillwell/GoRez/internal/request_service/interfaces"

	session "github.com/JackStillwell/GoRez/internal/session_service"
	sessionI "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
)

// has a struct that contains the other services
// manages limits and sessions
// is the recommended way for people to interact with the package

type svc struct {
	authSvc    authI.AuthService
	requestSvc requestI.RequestService
	sessionSvc sessionI.SessionService
}

type g struct {
	i.GodItemInfo
	i.PlayerInfo
	i.MatchInfo
}

func NewGorez() i.GoRez {
	s := svc{
		authSvc:    auth.NewAuthService(authM.Auth{}),
		requestSvc: request.NewRequestService(0),
		sessionSvc: session.NewSessionService(0, nil),
	}

	return &g{
		GodItemInfo: NewGodItemInfo(c.NewHiRezConstants(), s.requestSvc, s.authSvc, s.sessionSvc),
		PlayerInfo:  NewPlayerInfo(s.requestSvc, s.authSvc, s.sessionSvc),
		MatchInfo:   NewMatchInfo(s.requestSvc, s.authSvc, s.sessionSvc),
	}
}

func (r *g) Init() {

}
