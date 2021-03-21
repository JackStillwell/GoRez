package gorez

import (
	requestService "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"
)

type apiUtil struct {
	rS requestService.RequestService
}

func NewAPIUtil(rS requestService.RequestService) i.APIUtil {
	return &apiUtil{
		rS: rS,
	}
}

func (a *apiUtil) CreateSession() *m.Session {
	return nil
}

func (a *apiUtil) TestSession(s *m.Session) string {
	return ""
}

func (a *apiUtil) GetDataUsed(s *m.Session) *m.UsageInfo {
	return nil
}
