package gorezinternal

import "testing"
import mocker "github.com/JackStillwell/GoRez/test"
import "github.com/golang/mock/gomock"


func TestGetSignature(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockObj := mocker.NewMockHTTPGetter(mockCtrl)
	mockRM := RequestManager{}.mock(mockObj)

	want := "d1052b8673d7f6510e494721c129d306"
	
	if got := mockRM.getSignature("mockEndpoint", "00000000000000"); got != want {
		t.Errorf("getSignature() = %q, want %q", got, want)
	}
}