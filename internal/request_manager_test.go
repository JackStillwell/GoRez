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

	if got := requestManager.getSignature("mockEndpoint", "00000000000000"); got != want {
		t.Errorf("getSignature() = %q, want %q", got, want)
	}
}

func TestGetTimestamp(t *testing.T) {
	want := "00010101000000"
	dummyTime, _ := time.Parse("01 02 15 04 05 2006", "01 01 00 00 00 0001")

	if got := getTimestamp(dummyTime); got != want {
		t.Errorf("getTimestamp() = %q, want %q", got, want)
	}
}

func TestEndpointRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRequester := mocker.NewMockHTTPGetter(mockCtrl)

	urlStart := "mockURLBase/mockEndpointjson/mockDevID/"
	dummySignature := "f45bc3f241c136c2e808b8eb0f3891e9"
	dummySession := "/00000000000000/00010101000000"
	url := urlStart + dummySignature + dummySession

	mockRequester.EXPECT().Get(url).Return([]byte("success!"), nil)

	requestManager := RequestManager{}.mock(mockRequester)
	dummyTime, _ := time.Parse("01 02 15 04 05 2006", "01 01 00 00 00 0001")

	want := "success!"
	got, err := requestManager.EndpointRequest("mockEndpoint", "00000000000000", "", dummyTime)

	if string(got) != want || err != nil {
		t.Errorf("EndpointRequest() = %q, want %q, err %q", got, want, err)
	}
}

func TestEndpointRequestWithArgs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRequester := mocker.NewMockHTTPGetter(mockCtrl)

	urlStart := "mockURLBase/mockEndpointjson/mockDevID/"
	dummySignature := "f45bc3f241c136c2e808b8eb0f3891e9"
	dummySession := "/00000000000000/00010101000000/mockArgs"
	url := urlStart + dummySignature + dummySession

	mockRequester.EXPECT().Get(url).Return([]byte("success!"), nil)

	requestManager := RequestManager{}.mock(mockRequester)
	dummyTime, _ := time.Parse("01 02 15 04 05 2006", "01 01 00 00 00 0001")

	want := "success!"
	got, err := requestManager.EndpointRequest("mockEndpoint", "00000000000000", "mockArgs", dummyTime)

	if string(got) != want || err != nil {
		t.Errorf("EndpointRequest() = %q, want %q, err %q", got, want, err)
	}
}

func TestCreaseSessionRequest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRequester := mocker.NewMockHTTPGetter(mockCtrl)

	urlStart := "mockURLBase/createsessionjson/mockDevID/"
	dummySignature := "d15060234731b75d80ed365a8692ecb0"
	dummyTimestamp := "/00010101000000"
	url := urlStart + dummySignature + dummyTimestamp

	mockRequester.EXPECT().Get(url).Return([]byte("success!"), nil)

	requestManager := RequestManager{}.mock(mockRequester)
	dummyTime, _ := time.Parse("01 02 15 04 05 2006", "01 01 00 00 00 0001")

	want := "success!"
	got, err := requestManager.CreateSessionRequest(dummyTime)

	if string(got) != want || err != nil {
		t.Errorf("CreateSessionRequest() = %q, want %q, err %q", got, want, err)
	}
}
