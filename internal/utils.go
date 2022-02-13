package internal

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/onsi/ginkgo/v2"
	"github.com/pkg/errors"

	rSI "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	rSM "github.com/JackStillwell/GoRez/internal/request_service/models"
	sSI "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
	sSM "github.com/JackStillwell/GoRez/internal/session_service/models"
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

func BulkAsyncSessionRequest(rqstSvc rSI.RequestService, sesnSvc sSI.SessionService,
	requestBuilders []func(*sSM.Session) *rSM.Request) ([][]byte, []error) {
	numRequests := len(requestBuilders)
	responseChan := make(chan *rSM.RequestResponse, numRequests)
	uIDSessionMap := make(map[*uuid.UUID]*sSM.Session, numRequests)

	// NOTE: this is async so the reservation and release of sessions is possible, but the func
	// return depends upon responses being completed.
	go func() {
		// WARN: this is required for testing recovery. Uncomment if debugging test failures.
		defer ginkgo.GinkgoRecover()

		for _, rB := range requestBuilders {
			sessChan := make(chan *sSM.Session, 1)
			sesnSvc.ReserveSession(1, sessChan)

			s := <-sessChan // NOTE: will wait here until session recieved

			// Request constructed with session here
			r := rB(s)

			uID := uuid.New()
			r.Id = &uID
			uIDSessionMap[&uID] = s

			rqstSvc.MakeRequest(r)
			rqstSvc.GetResponse(&uID, responseChan)
		}
	}()

	responses := make([][]byte, 0, numRequests)
	errs := make([]error, 0, numRequests)
	for i := 0; i < numRequests; i++ {
		resp := <-responseChan
		if resp.Err != nil {
			if strings.Contains(resp.Err.Error(), "session") {
				sesnSvc.BadSession([]*sSM.Session{uIDSessionMap[resp.Id]})
			} else {
				sesnSvc.ReleaseSession([]*sSM.Session{uIDSessionMap[resp.Id]})
			}
			errs = append(errs, errors.Wrap(resp.Err, "request"))
			continue
		}

		sesnSvc.ReleaseSession([]*sSM.Session{uIDSessionMap[resp.Id]})

		responses = append(responses, resp.Resp)
	}

	return responses, errs
}
