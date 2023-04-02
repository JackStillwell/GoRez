package gorez

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	c "github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"

	auth "github.com/JackStillwell/GoRez/internal/auth"
	authI "github.com/JackStillwell/GoRez/internal/auth/interfaces"
	authM "github.com/JackStillwell/GoRez/internal/auth/models"

	request "github.com/JackStillwell/GoRez/internal/request"
	requestI "github.com/JackStillwell/GoRez/internal/request/interfaces"

	session "github.com/JackStillwell/GoRez/internal/session"
	sessionI "github.com/JackStillwell/GoRez/internal/session/interfaces"
	sessionM "github.com/JackStillwell/GoRez/internal/session/models"
)

// has a struct that contains the other services
// manages limits and sessions
// is the recommended way for people to interact with the package

type svc struct {
	AuthSvc    authI.Service
	RequestSvc requestI.Service
	SessionSvc sessionI.Service
}

type g struct {
	*svc
	i.APIUtil
	i.GodItemInfo
	i.PlayerInfo
	i.MatchInfo
}

func NewGorez(auth_path string) (i.GoRez, error) {
	contents, err := os.ReadFile(auth_path)
	if err != nil {
		return nil, errors.Wrap(err, "reading file")
	}
	lines := strings.Split(string(contents), "\n")
	if len(lines) < 2 || len(lines) > 3 {
		return nil, errors.New("auth file must contain two lines, the first being your dev id and " +
			"the second being your dev key")
	}

	s := &svc{
		AuthSvc:    auth.NewService(authM.Auth{ID: lines[0], Key: lines[1]}),
		RequestSvc: request.NewService(1),
		SessionSvc: session.NewService(1, nil),
	}

	util := NewGorezUtil(s.AuthSvc, s.RequestSvc, s.SessionSvc)

	return &g{
		svc:         s,
		APIUtil:     NewAPIUtil(c.NewHiRezConstants(), s.AuthSvc, s.RequestSvc, s.SessionSvc),
		GodItemInfo: NewGodItemInfo(c.NewHiRezConstants(), util),
		PlayerInfo:  NewPlayerInfo(s.RequestSvc, s.AuthSvc, s.SessionSvc),
		MatchInfo:   NewMatchInfo(s.RequestSvc, s.AuthSvc, s.SessionSvc),
	}, nil
}

func (gr *g) Init() error {
	log.Println("looking for sessions.json")
	if _, err := os.Stat("sessions.json"); err == nil {
		log.Println("sessions.json found")
		contents, err := os.ReadFile("sessions.json")
		if err != nil {
			return errors.Wrap(err, "reading sessions.txt")
		}

		var existingSessions []*sessionM.Session
		err = json.Unmarshal(contents, &existingSessions)
		if err != nil {
			return errors.Wrap(err, "unmarshaling sessions.txt")
		}
		gr.SessionSvc.ReleaseSession(existingSessions)
	} else if errors.Is(err, os.ErrNotExist) {
		log.Println("sessions.json not found - creating session")
		sessions, errs := gr.APIUtil.CreateSession(1)
		if errs[0] != nil {
			return errors.Wrap(errs[0], "creating session failed")
		}
		log.Println("session created")

		created, err := time.ParseInLocation("1/2/2006 3:04:05 PM", *sessions[0].Timestamp, time.UTC)
		if err != nil {
			return errors.Wrap(err, "parsing session timestamp")
		}
		internalSession := &sessionM.Session{
			Key:     *sessions[0].SessionID,
			Created: &created,
		}

		gr.SessionSvc.ReleaseSession([]*sessionM.Session{internalSession})
	} else {
		return fmt.Errorf("stat-ing sessions.json: %w", err)
	}

	return nil
}

func (gr *g) Shutdown() {
	// store the sessions here so they're not lost on each run

	jBytes, err := json.Marshal(gr.SessionSvc.GetAvailableSessions())
	if err != nil {
		log.Println("error marshaling items", err)
	}

	f, err := os.Create("sessions.json")
	if err != nil {
		log.Println("error writing sessions", err)
		return
	}
	defer f.Close()
	f.Write(jBytes)
}
