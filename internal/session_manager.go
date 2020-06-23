package gorezinternal

import (
	"encoding/json"
	"errors"
	"sync"
)

// SessionManager contains the information necessary for managing HiRez API Session
type SessionManager struct {
	requestManager  RequestManagement
	idleSessions    []string
	activeSessions  []string
	sessionsCreated uint16
	mux             sync.Mutex
}

func (t *SessionManager) mock(rm RequestManagement) {
	t.requestManager = rm
	t.idleSessions = make([]string, 0)
	t.activeSessions = make([]string, 0)
	t.sessionsCreated = 0
}

// Initialize returns a SessionManager constructed with accurate startup
//   values attempts to reuse old sessions are made and current limit info is
//   pulled down from the HiRez API
func (t *SessionManager) Initialize(initFile string, rm RequestManagement) {
	limitConsts := LimitConstants{}.New()

	// read existing sessions from file
	iS := make([]string, limitConsts.ConcurrentSessions)

	apiConsts := APIConstants{}.New()

	t.requestManager = rm
	t.idleSessions = iS
	t.activeSessions = make([]string, limitConsts.ConcurrentSessions)
	t.sessionsCreated = 0

	sessionid, err := t.GetSession()

	// check the err
	if err != nil {
		panic("Could not acquire session to determine current usage")
	}

	// read current sessionsCreated from the API
	rm.EndpointRequest(
		apiConsts.TestSession,
		sessionid,
		"",
	)

	// parse from json obj

	// populate the sm
	t.sessionsCreated = 0
}

// Save stores the sessions currently created in a file specified
//   NOTE: THIS SHOULD BE DEFFERED ON CREATION TO ENSURE SESSIONS ARE NOT LOST
func (t *SessionManager) Save(saveFile string) {
	// saves stuff
}

// Session contains the information returned from a createsession request
type Session struct {
	retMsg    string
	sessionID string
	timestamp string
}

// ParseJSONToSession parses a json []byte to a Session struct
func ParseJSONToSession(jsonString []byte) (Session, error) {
	var rawMap map[string]interface{}
	jsonErr := json.Unmarshal(jsonString, &rawMap)
	if jsonErr != nil {
		return Session{}, jsonErr
	}

	retMsg, _ := rawMap["ret_msg"].(string)

	// NOTE: need to check retmsg to see if I should attempt
	//         to parse the other fields

	sessionID, _ := rawMap["session_id"].(string)
	timestamp, _ := rawMap["timestamp"].(string)

	session := Session{
		retMsg:    retMsg,
		sessionID: sessionID,
		timestamp: timestamp,
	}

	return session, nil
}

// GetSession provides a valid session id or returns an error
func (t *SessionManager) GetSession() (string, error) {
	t.mux.Lock()

	var sessionid string
	if len(t.idleSessions) > 0 {
		sessionid = t.idleSessions[0]
		t.idleSessions = t.idleSessions[1:]
		t.activeSessions = append(t.activeSessions, sessionid)

		apiConsts := APIConstants{}.New()

		// test the session to make sure its still valid
		resp, err := t.requestManager.EndpointRequest(
			apiConsts.TestSession,
			sessionid,
			"",
		)

		if err != nil {
			return "", err
		}

		// parse resp
		if string(resp) != "write a parser for the testsession reply" {
			return "", errors.New("not implemented")
		}

	} else {
		limitConsts := LimitConstants{}.New()

		if t.sessionsCreated == limitConsts.SessionsPerDay {
			return "", errors.New("Sessions Created per Day Limit Reached")
		}

		body, requestErr := t.requestManager.CreateSessionRequest()

		if requestErr != nil {
			return "", requestErr
		}

		session, jsonErr := ParseJSONToSession(body)
		if jsonErr != nil {
			return "", jsonErr
		}

		sessionid = session.sessionID
		t.activeSessions = append(t.activeSessions, sessionid)
		t.sessionsCreated++
	}

	t.mux.Unlock()
	return sessionid, nil
}

// ReturnSession places a sessionid in the pool for use
func (t *SessionManager) ReturnSession(sessionID string) error {
	t.mux.Lock()

	foundSession := false
	var i int
	var v string
	for i, v = range t.activeSessions {
		if v == sessionID {
			foundSession = true
			break
		}
	}

	if !foundSession {
		return errors.New("Did not find session in active sessions")
	}

	firstPart := t.activeSessions[:i]
	secondPart := t.activeSessions[i+1:]
	t.activeSessions = append(firstPart, secondPart...)
	t.idleSessions = append(t.idleSessions, v)

	t.mux.Unlock()
	return nil
}
