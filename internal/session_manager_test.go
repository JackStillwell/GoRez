package internal_gorez

import (
	"testing"

	mocker "github.com/JackStillwell/GoRez/test"
	"github.com/golang/mock/gomock"
)

func TestParseJSONToSession(t *testing.T) {
	json := []byte(`{
    "ret_msg": "Approved",
    "session_id": "dummy_id",
    "timestamp": "3/29/2020 3:12:06 PM"
}`)

	want := Session{
		retMsg:    "Approved",
		sessionID: "dummy_id",
		timestamp: "3/29/2020 3:12:06 PM",
	}
	if got, err := ParseJSONToSession(json); got != want || err != nil {
		t.Errorf("ParseJSONToSession() = %q, want %q err %q", got, want, err.Error())
	}
}

func TestGetSession(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockObj := mocker.NewMockHttpGetter(mockCtrl)
	mockObj.EXPECT().Get("http://api.smitegame.com/smiteapi.svc/json/dummy_id/dummy_key").Return([]byte(`{
    "ret_msg": "Approved",
    "session_id": "dummy_id",
    "timestamp": "3/29/2020 3:12:06 PM"
}`), nil)

	want := "dummy_id"

	if got, err := GetSession("dummy_id", "dummy_key", mockObj); got != want || err != nil {
		t.Errorf("GetSession() = %q, want %q err %q", got, want, err.Error())
	}
}
