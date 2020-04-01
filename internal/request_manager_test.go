package gorezinternal

import (
	"testing"
	"time"

	mocker "github.com/JackStillwell/GoRez/test"
	"github.com/golang/mock/gomock"
)

func TestGetSignature(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockObj := mocker.NewMockHTTPGetter(mockCtrl)
	requestManager := RequestManager{}.mock(mockObj)

	want := "d1052b8673d7f6510e494721c129d306"

	if got := requestManager.GetSignature("mockEndpoint", "00000000000000"); got != want {
		t.Errorf("getSignature() = %q, want %q", got, want)
	}
}

func TestGetTimestamp(t *testing.T) {
	want := "00010101000000"
	dummyTime, _ := time.Parse("01 02 15 04 05 2006", "01 01 00 00 00 0001")

	if got := GetTimestamp(dummyTime); got != want {
		t.Errorf("getTimestamp() = %q, want %q", got, want)
	}
}

func TestEndpointRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRequester := mocker.NewMockHTTPGetter(mockCtrl)

	url := "mockURLBase/mockEndpointjson/mockDevID/f45bc3f241c136c2e808b8eb0f3891e9/00000000000000/00010101000000/"
	mockRequester.EXPECT().Get(url).Return([]byte("success!"), nil)

	requestManager := RequestManager{}.mock(mockRequester)
	dummyTime, _ := time.Parse("01 02 15 04 05 2006", "01 01 00 00 00 0001")

	want := "success!"

	if got, err := requestManager.EndpointRequest("mockEndpoint", "00000000000000", "", dummyTime); string(got) != want || err != nil {
		t.Errorf("getSignature() = %q, want %q, err %q", got, want, err)
	}
}
