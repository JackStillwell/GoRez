package gorez_test

import (
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gorez "github.com/JackStillwell/GoRez/pkg"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"

	requestMocks "github.com/JackStillwell/GoRez/internal/request_service/mocks"
	requestM "github.com/JackStillwell/GoRez/internal/request_service/models"

	sessionMocks "github.com/JackStillwell/GoRez/internal/session_service/mocks"
	sessionM "github.com/JackStillwell/GoRez/internal/session_service/models"

	authMocks "github.com/JackStillwell/GoRez/internal/auth_service/mocks"
)

var _ = Describe("GorezUtil", func() {
	var (
		ctrl *gomock.Controller

		rqstSvc *requestMocks.MockRequestService
		sesnSvc *sessionMocks.MockSessionService
		authSvc *authMocks.MockAuthService

		target i.GorezUtil
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		rqstSvc = requestMocks.NewMockRequestService(ctrl)
		sesnSvc = sessionMocks.NewMockSessionService(ctrl)
		authSvc = authMocks.NewMockAuthService(ctrl)

		target = gorez.NewGorezUtil(authSvc, rqstSvc, sesnSvc)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("BulkAsyncSessionRequest", func() {
		It("should return each response", func() {
			sess := &sessionM.Session{
				Key: "fake",
			}
			sessChan := make(chan *sessionM.Session, 1)
			sesnSvc.EXPECT().ReserveSession(1, gomock.AssignableToTypeOf(sessChan)).
				Do(func(_ int, sessChan chan *sessionM.Session) {
					sessChan <- sess
				})

			requestBuilder := func(s *sessionM.Session) *requestM.Request {
				return &requestM.Request{
					JITArgs: []interface{}{s},
					JITBuild: func(...any) (string, error) {
						return "resp", nil
					},
				}
			}

			request := requestBuilder(sess)
			rqstSvc.EXPECT().MakeRequest(gomock.AssignableToTypeOf(request)).Do(func(r *requestM.Request) {
				request.Id = r.Id
				Expect(r.JITArgs).To(Equal(request.JITArgs))
				Expect(r.JITBuild).To(BeAssignableToTypeOf(request.JITBuild))
			})

			respChan := make(chan *requestM.RequestResponse, 1)
			rqstSvc.EXPECT().GetResponse(
				gomock.AssignableToTypeOf(&uuid.UUID{}),
				gomock.AssignableToTypeOf(respChan),
			).Do(func(uID *uuid.UUID, respChan chan *requestM.RequestResponse) {
				respChan <- &requestM.RequestResponse{
					Id:   uID,
					Err:  nil,
					Resp: []byte("stuff"),
				}
			})

			sesnSvc.EXPECT().ReleaseSession(gomock.AssignableToTypeOf([]*sessionM.Session{})).Do(
				func(sessions []*sessionM.Session) {
					Expect(sessions).To(ConsistOf(sess))
				},
			)

			responsesChan := make(chan [][]byte, 1)
			errorsChan := make(chan []error, 1)
			go func() {
				defer GinkgoRecover()

				r, e := target.BulkAsyncSessionRequest(
					[]func(*sessionM.Session) *requestM.Request{requestBuilder},
				)
				responsesChan <- r
				errorsChan <- e
			}()

			Eventually(responsesChan).Should(Receive(ConsistOf([][]byte{[]byte("stuff")})))
			Eventually(errorsChan).Should(Receive(&[]error{nil}))
		})
	})
})
