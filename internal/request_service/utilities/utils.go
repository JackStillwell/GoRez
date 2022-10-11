package utils

import (
	"errors"
	"fmt"
	"time"
)

type JITFunc func(args []any) (string, error)

// Utilities should have no dependencies in the project, so they can be cross-imported without
// worrying about dependency cycles.

/*
JITBase takes the following args:
 1. baseURL string
 2. devID string
 3. endpoint string
 4. session string ["" if none]
 5. timeStamp func(time.Time) string
 6. signature func(endpoint, timeStamp string) string
 7. endpointArgs string ["" if none]
*/
func JITBase(args ...any) (string, error) {

	if len(args) != 7 {
		return "", errors.New("incorrect number of arguments passed")
	}

	baseURL, ok := args[0].(string)
	if !ok {
		return "", errors.New("could not coerce first arg to string")
	}

	devID, ok := args[1].(string)
	if !ok {
		return "", errors.New("could not coerce second arg to string")
	}

	t := time.Now().UTC()

	endpoint, ok := args[2].(string)
	if !ok {
		return "", errors.New("could not coerce third arg to string")
	}

	session, ok := args[3].(string)
	if !ok {
		return "", errors.New("could not coerce fourth arg to string")
	}

	tS, ok := args[4].(func(time.Time) string)
	if !ok {
		return "", errors.New("could not coerce fifth arg to func(time.Time) string")
	}

	s, ok := args[5].(func(endpoint, timeStamp string) string)
	if !ok {
		return "", errors.New("could not coerce sixth arg to func(endpoint, timeStamp string) string")
	}

	endpointArgs, ok := args[6].(string)
	if !ok {
		return "", errors.New("could not coerce seventh arg to string")
	}

	timeStamp := tS(t)
	signature := s(endpoint, timeStamp)

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
