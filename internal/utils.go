package gorezinternal

import (
	"io/ioutil"
	"net/http"
)

// DefaultGetter is the default HTTPGetter implementation
type DefaultGetter struct{}

// Get retrieves a byte array from a URL
func (t DefaultGetter) Get(url string) ([]byte, error) {
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
