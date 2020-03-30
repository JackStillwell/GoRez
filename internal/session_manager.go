package gorezinternal

import (
	"encoding/json"
	"fmt"
)

// Session has no comment
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
func GetSession(
	baseURL,
	returnDataType,
	devID,
	devKey string,
	httpGet HTTPGetter) (string, error) {

	request := fmt.Sprintf(
		"%s/%s/%s/%s",
		baseURL,
		returnDataType,
		devID,
		devKey,
	)
	body, getterErr := httpGet.Get(request)

	if getterErr != nil {
		return "HttpGetter error", getterErr
	}

	session, jsonErr := ParseJSONToSession(body)
	if jsonErr != nil {
		return "json Unmarshal error", jsonErr
	}

	return session.sessionID, nil
}
