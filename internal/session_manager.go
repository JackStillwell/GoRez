package gorezinternal

import (
	"encoding/json"
	"errors"
	"sync"
)

// SessionManager contains the information necessary for managing HiRez API Session
type SessionManager struct {
	requestManager  RequestManager
	idleSessions    []string
	activeSessions  []string
	sessionsCreated uint8
	mux             sync.Mutex
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

	if len(t.idleSessions) > 0 {
		toreturn := t.idleSessions[0]
		t.idleSessions = t.idleSessions[1:]
		t.activeSessions = append(t.activeSessions, toreturn)
		return toreturn, nil
	}
	body, requestErr := t.requestManager.CreateSessionRequest()

	if requestErr != nil {
		return "", requestErr
	}

	session, jsonErr := ParseJSONToSession(body)
	if jsonErr != nil {
		return "", jsonErr
	}

	t.mux.Unlock()
	return session.sessionID, nil
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
