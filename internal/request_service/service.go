package request_service

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	i "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	m "github.com/JackStillwell/GoRez/internal/request_service/models"
	"github.com/pkg/errors"
)

type requestService struct {
	i.Requester
	i.RequestManager
}

type requester struct {
	http i.HTTPGet
}

type requestManager struct {
	r            i.Requester
	requestChan  chan *m.Request
	responseChan chan *m.RequestResponse
	workerKill   []chan bool
}

type httpGetter struct{}

func (*httpGetter) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

func NewTestRequester(http i.HTTPGet) i.Requester {
	return &requester{http}
}

func NewTestRequestService(capacity int, r i.Requester) i.RequestService {
	requests := make(chan *m.Request, capacity)
	responses := make(chan *m.RequestResponse, capacity)

	rM := &requestManager{
		r:            r,
		requestChan:  requests,
		responseChan: responses,
	}

	wKs := make([]chan bool, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		wKs[i] = make(chan bool)
		go requestServiceRoutine(rM, wKs[i])
	}

	rM.workerKill = wKs

	rS := &requestService{r, rM}

	return rS
}

func NewRequestService(capacity int) i.RequestService {
	return NewTestRequestService(capacity, NewTestRequester(&httpGetter{}))
}

func (r *requester) Request(rqst *m.Request) (rr *m.RequestResponse) {
	rr = &m.RequestResponse{
		Id:   rqst.Id,
		Resp: nil,
		Err:  nil,
	}

	requestURL, err := rqst.JITBuild(rqst.JITArgs)
	if err != nil {
		rr.Err = errors.Wrap(err, "building requesturl")
		return
	}

	resp, err := r.http.Get(requestURL)
	if err != nil {
		rr.Err = errors.Wrap(err, "getting response")
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		rr.Err = errors.Wrap(err, "reading body")
		return
	}

	rr.Resp = body
	return
}

func (rM *requestManager) MakeRequest(r *m.Request) {
	rM.requestChan <- r
}

func (rM *requestManager) GetResponse() (toRet *m.RequestResponse) {
	return <-rM.responseChan
}

func (rM *requestManager) Close() {
	for _, c := range rM.workerKill {
		c <- true
	}
}

func requestServiceRoutine(rM *requestManager, killChan chan bool) {
	kill := false
	for !kill {
		select {
		case rqst := <-rM.requestChan:
			rM.responseChan <- rM.r.Request(rqst)
		case <-killChan:
			kill = true
		}
	}
}

/* JITBase takes the following args:
   1) baseURL string
   2) devID string
   3) endpoint string
   4) session string
   5) timeStamp func(time.Time) string
   6) signature func(endpoint, timeStamp string) string
   7) endpointArgs string
*/
func (*requester) JITBase(args []interface{}) (string, error) {
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
