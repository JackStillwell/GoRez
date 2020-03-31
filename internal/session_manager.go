package gorezinternal

import (
	"encoding/json"
	"fmt"
)

// SessionManager contains the information necessary for managing HiRez API Session
type SessionManager struct {
	idleSessions    []string
	activeSessions  []string
	sessionsCreated uint8
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

// GetSession gets a new session
func GetSession(api APIBase) (string, error) {

	request := fmt.Sprintf(
		"%s/%s/%s/%s",
		api.baseURL,
		api.returnDataType,
		api.devID,
		api.devKey,
	)
	body, getterErr := api.httpGet.Get(request)

	if getterErr != nil {
		return "", getterErr
	}

	session, jsonErr := ParseJSONToSession(body)
	if jsonErr != nil {
		return "", jsonErr
	}

	return session.sessionID, nil
}
