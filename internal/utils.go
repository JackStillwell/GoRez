package gorez_internal

import (
	"io/ioutil"
	"net/http"
)

// defaultGetter is the default httpGetter implementation
type defaultGetter struct{}

// Get retrieves a byte array from a URL
func (t defaultGetter) get(url string) ([]byte, error) {
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
