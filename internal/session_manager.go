package gorez

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Session has no comment
type Session struct {
	retMsg    string
	sessionID string
	timestamp string
}

// Get retrieves a byte array from a URL
func Get(url string) ([]byte, error) {
	resp, getErr := http.Get(url)

	if getErr != nil {
		return nil, getErr
	}

	defer resp.Body.Close()
	body, readErr := ioutil.ReadAll(resp.Body)

	if readErr != nil {
		return nil, readErr
	}

	return body, nil
}

// ParseJSONToSession parses a json []byte to a Session struct
func ParseJSONToSession(jsonString []byte) (Session, error) {
	var rawMap map[string]interface{}
	jsonErr := json.Unmarshal(jsonString, &rawMap)
	if jsonErr != nil {
		return Session{}, jsonErr
	}

	retMsg, _ := rawMap["ret_msg"].(string)
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
func GetSession(id, key string, getter HTTPGetter) (string, error) {
	request := fmt.Sprintf("http://api.smitegame.com/smiteapi.svc/json/%s/%s", id, key)
	body, getterErr := getter.Get(request)

	if getterErr != nil {
		return "HttpGetter error", getterErr
	}

	session, jsonErr := ParseJSONToSession(body)
	if jsonErr != nil {
		return "json Unmarshal error", jsonErr
	}

	return session.sessionID, nil
}
