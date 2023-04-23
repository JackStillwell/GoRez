package gorez

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

const NUM_SESSIONS = 40

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
	sessionCache i.SessionCache
}

type localSessionCache struct{}

func (localSessionCache) ReadSessions() ([]*sessionM.Session, error) {
	log.Println("looking for sessions.json")
	if _, err := os.Stat("sessions.json"); err == nil {
		log.Println("sessions.json found")
		contents, err := os.ReadFile("sessions.json")
		if err != nil {
			return nil, fmt.Errorf("reading sessions.txt: %w", err)
		}

		var existingSessions []*sessionM.Session
		err = json.Unmarshal(contents, &existingSessions)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling sessions.txt: %w", err)
		}

		return existingSessions, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return []*sessionM.Session{}, nil
	} else {
		return nil, fmt.Errorf("stat-ing sessions.json: %w", err)
	}
}

func (localSessionCache) SaveSessions(sessions []*sessionM.Session) error {
	jBytes, err := json.Marshal(sessions)
	if err != nil {
		return fmt.Errorf("marshaling items: %w", err)
	}

	f, err := os.Create("sessions.json")
	if err != nil {
		return fmt.Errorf("creating session file: %w", err)
	}
	defer f.Close()
	if _, err := f.Write(jBytes); err != nil {
		return fmt.Errorf("writing sessions: %w", err)
	}

	return nil
}

func NewGorez(auth_path string, sessionCache i.SessionCache) (i.GoRez, error) {
	contents, err := os.ReadFile(auth_path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	lines := strings.Split(string(contents), "\n")
	if len(lines) < 2 || len(lines) > 3 {
		return nil, fmt.Errorf("auth file must contain two lines, the first being your dev id and " +
			"the second being your dev key")
	}

	s := &svc{
		AuthSvc:    auth.NewService(authM.Auth{ID: lines[0], Key: lines[1]}),
		RequestSvc: request.NewService(NUM_SESSIONS),
		SessionSvc: session.NewService(NUM_SESSIONS, nil),
	}

	util := NewGorezUtil(s.AuthSvc, s.RequestSvc, s.SessionSvc)

	if sessionCache == nil {
		sessionCache = localSessionCache{}
	}

	return &g{
		svc:          s,
		APIUtil:      NewAPIUtil(c.NewHiRezConstants(), s.AuthSvc, s.RequestSvc, s.SessionSvc),
		GodItemInfo:  NewGodItemInfo(c.NewHiRezConstants(), util),
		PlayerInfo:   NewPlayerInfo(s.RequestSvc, s.AuthSvc, s.SessionSvc),
		MatchInfo:    NewMatchInfo(s.RequestSvc, s.AuthSvc, s.SessionSvc),
		sessionCache: sessionCache,
	}, nil
}

func (gr *g) createSessions(numSessions int) error {
	log.Printf("creating %d sessions\n", numSessions)
	sessions, errs := gr.APIUtil.CreateSession(numSessions)
	errCount := 0
	for i, e := range errs {
		if e != nil {
			log.Printf("error creating session %d: %s\n", i, e.Error())
			errCount++
		}
	}
	if errCount == numSessions {
		return fmt.Errorf("all session creations errored")
	}
	log.Println("sessions created")

	sessionObjs := make([]*sessionM.Session, 0, numSessions)
	for i, session := range sessions {
		if session != nil {
			created, err := time.ParseInLocation("1/2/2006 3:04:05 PM", *session.Timestamp, time.UTC)
			if err != nil {
				log.Printf("parsing session timestamp for session %d: %s\n", i, err.Error())
				continue
			}
			sessionObjs = append(sessionObjs, &sessionM.Session{
				Key:     *session.SessionID,
				Created: &created,
			})
		}
	}

	gr.SessionSvc.ReleaseSession(sessionObjs)
	return nil
}

func (gr *g) Init() error {

	// get stored sessions
	existingSessions, err := gr.sessionCache.ReadSessions()
	if err != nil {
		log.Printf("error reading sessions: %s", err.Error())
	}

	// test sessions
	sessionKeys := make([]string, 0, len(existingSessions))
	for _, eS := range existingSessions {
		sessionKeys = append(sessionKeys, eS.Key)
	}

	validSessions := make([]*sessionM.Session, 0, len(existingSessions))
	responses, errs := gr.APIUtil.TestSession(sessionKeys)
	for i, resp := range responses {
		if resp != nil {
			if !strings.Contains(*resp, "Invalid session id") {
				validSessions = append(validSessions, existingSessions[i])
			}
		} else {
			log.Printf("error testing session %s: %s", existingSessions[i].Key, errs[i].Error())
		}
	}

	gr.SessionSvc.ReleaseSession(validSessions)
	if err := gr.createSessions(NUM_SESSIONS - len(validSessions)); err != nil {
		return fmt.Errorf("creating sessions: %w", err)
	}

	return nil
}

func (gr *g) Shutdown() {
	// store the sessions here so they're not lost on each run

	if err := gr.sessionCache.SaveSessions(gr.SessionSvc.GetAvailableSessions()); err != nil {
		log.Printf("saving sessions: %s", err.Error())
	}
}
