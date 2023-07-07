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
	m "github.com/JackStillwell/GoRez/pkg/models"
	"go.uber.org/zap"

	auth "github.com/JackStillwell/GoRez/internal/auth"
	authI "github.com/JackStillwell/GoRez/internal/auth/interfaces"
	authM "github.com/JackStillwell/GoRez/internal/auth/models"
	"github.com/JackStillwell/GoRez/internal/base"

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
	b          base.Service
}

type g struct {
	*svc
	i.APIUtil
	i.GodItemInfo
	i.PlayerInfo
	i.MatchInfo
	sessionCache i.SessionCache
	numSessions  int
}

func hirezToInternalSessions(internalSessions []*m.Session) []*sessionM.Session {
	sessionObjs := make([]*sessionM.Session, 0, len(internalSessions))
	for i, session := range internalSessions {
		if session != nil {
			created, err := time.ParseInLocation("1/2/2006 3:04:05 PM", *session.Timestamp, time.UTC)
			if err != nil {
				log.Printf("parsing session timestamp for session %d: %s", i, err.Error())
				continue
			}
			sessionObjs = append(sessionObjs, &sessionM.Session{
				Key:     *session.SessionID,
				Created: &created,
			})
		}
	}

	return sessionObjs
}

func internalToHirezSessions(hirezSessions []*sessionM.Session) []*m.Session {
	sessionObjs := make([]*m.Session, 0, len(hirezSessions))
	for _, session := range hirezSessions {
		if session != nil {
			timestamp := session.Created.Format("1/2/2006 3:04:05 PM")
			sessionObjs = append(sessionObjs, &m.Session{
				SessionID: &session.Key,
				Timestamp: &timestamp,
			})
		}
	}

	return sessionObjs
}

type localSessionCache struct{}

func (localSessionCache) ReadSessions() ([]*m.Session, error) {
	log.Println("looking for sessions.json")
	if _, err := os.Stat("sessions.json"); err == nil {
		log.Println("sessions.json found")
		contents, err := os.ReadFile("sessions.json")
		if err != nil {
			return nil, fmt.Errorf("reading sessions.txt: %w", err)
		}

		var existingSessions []*m.Session
		err = json.Unmarshal(contents, &existingSessions)
		if err != nil {
			return nil, fmt.Errorf("unmarshaling sessions.txt: %w", err)
		}

		return existingSessions, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return []*m.Session{}, nil
	} else {
		return nil, fmt.Errorf("stat-ing sessions.json: %w", err)
	}
}

func (localSessionCache) SaveSessions(sessions []*m.Session) error {
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

func NewGorez(auth_path string, sessionCache i.SessionCache, numSessions int, logger *zap.Logger) (i.GoRez, error) {
	contents, err := os.ReadFile(auth_path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	lines := strings.Split(string(contents), "\n")
	if len(lines) < 2 || len(lines) > 3 {
		return nil, fmt.Errorf("auth file must contain two lines, the first being your dev id and " +
			"the second being your dev key")
	}

	b := base.NewService(logger)
	s := &svc{
		AuthSvc:    auth.NewService(authM.Auth{ID: lines[0], Key: lines[1]}, b),
		RequestSvc: request.NewService(numSessions, b),
		SessionSvc: session.NewService(numSessions, nil, b),
		b:          b,
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
		numSessions:  numSessions,
	}, nil
}

func (gr *g) createSessions(numSessions int) error {
	if numSessions <= 0 {
		return nil
	}
	log := gr.b.GetLogger()
	log.Info("creating sessions", zap.Int("numSessions", numSessions))
	sessions, errs := gr.APIUtil.CreateSession(numSessions)
	errCount := 0
	for i, e := range errs {
		if e != nil {
			log.Error("failed creating session", zap.Error(e), zap.Int("sessionNum", i))
			errCount++
		}
	}
	if errCount == numSessions {
		return fmt.Errorf("all session creations errored")
	}
	log.Info("sessions created")

	sessionObjs := hirezToInternalSessions(sessions)

	gr.SessionSvc.ReleaseSession(sessionObjs)
	return nil
}

func (gr *g) Init() error {
	log := gr.b.GetLogger()

	// get stored sessions
	existingSessions, err := gr.sessionCache.ReadSessions()
	if err != nil {
		log.Error("failed reading sessions", zap.Error(err))
	}

	// test sessions
	sessionKeys := make([]string, 0, len(existingSessions))
	for _, eS := range existingSessions {
		sessionKeys = append(sessionKeys, *eS.SessionID)
	}

	validSessions := make([]*sessionM.Session, 0, len(existingSessions))
	responses, errs := gr.APIUtil.TestSession(sessionKeys)
	for i, resp := range responses {
		if resp != nil {
			if !strings.Contains(*resp, "Invalid session id") {
				created, err := time.ParseInLocation("1/2/2006 3:04:05 PM",
					*existingSessions[i].Timestamp, time.UTC)
				if err != nil {
					log.Error("failure parsing session timestamp", zap.Error(err),
						zap.String("timestamp", *existingSessions[i].Timestamp))
					continue
				}
				validSessions = append(validSessions, &sessionM.Session{
					Key:     *existingSessions[i].SessionID,
					Created: &created,
				})
			}
		} else {
			log.Error("failure testing session",
				zap.String("sessionID", *existingSessions[i].SessionID), zap.Error(errs[i]))
		}
	}

	gr.SessionSvc.ReleaseSession(validSessions)
	numToCreate := gr.numSessions - len(validSessions)
	if numToCreate < 0 {
		numToCreate = 0
	}
	if err := gr.createSessions(numToCreate); err != nil {
		return fmt.Errorf("creating sessions: %w", err)
	}

	return nil
}

func (gr *g) Shutdown() {
	// store the sessions here so they're not lost on each run
	log := gr.b.GetLogger()

	if err := gr.sessionCache.SaveSessions(internalToHirezSessions(
		gr.SessionSvc.GetAvailableSessions(),
	)); err != nil {
		log.Error("failure saving sessions", zap.Error(err))
	}

	gr.b.ShutdownLogger()
}
