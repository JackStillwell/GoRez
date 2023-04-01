package gorez

import (
	"fmt"
	"time"
)

func HiRezJIT(baseURL, devID, endpoint, session string, timeStampFunc func(time.Time) string,
	signatureFunc func(endpoint, timeStamp string) string, endpointArgs string,
) func() (string, error) {
	toRet := func() (string, error) {
		t := time.Now().UTC()

		timeStamp := timeStampFunc(t)
		signature := signatureFunc(endpoint, timeStamp)

		if endpointArgs == "" && session == "" {
			return fmt.Sprintf(
				"%s/%s/%s/%s",
				baseURL,
				devID,
				signature,
				timeStamp,
			), nil
		} else if endpointArgs == "" {
			return fmt.Sprintf(
				"%s/%s/%s/%s/%s",
				baseURL,
				devID,
				signature,
				session,
				timeStamp,
			), nil
		} else {
			return fmt.Sprintf(
				"%s/%s/%s/%s/%s/%s",
				baseURL,
				devID,
				signature,
				session,
				timeStamp,
				endpointArgs,
			), nil
		}
	}

	return toRet
}
