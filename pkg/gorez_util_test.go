package gorez_test

import (
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gorez "github.com/JackStillwell/GoRez/pkg"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"

	requestMocks "github.com/JackStillwell/GoRez/internal/request/mocks"
	requestM "github.com/JackStillwell/GoRez/internal/request/models"

	sessionMocks "github.com/JackStillwell/GoRez/internal/session/mocks"
	sessionM "github.com/JackStillwell/GoRez/internal/session/models"

	authMocks "github.com/JackStillwell/GoRez/internal/auth/mocks"
)

var _ = Describe("GorezUtil", func() {
	var (
		ctrl *gomock.Controller

		rqstSvc *requestMocks.MockService
		sesnSvc *sessionMocks.MockService
		authSvc *authMocks.MockService

		target i.GorezUtil
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())

		rqstSvc = requestMocks.NewMockService(ctrl)
		sesnSvc = sessionMocks.NewMockService(ctrl)
		authSvc = authMocks.NewMockService(ctrl)

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
					JITFunc: func() (string, error) {
						return "resp", nil
					},
				}
			}

			request := requestBuilder(sess)
			rqstSvc.EXPECT().MakeRequest(gomock.AssignableToTypeOf(request)).Do(func(r *requestM.Request) {
				request.Id = r.Id
				url, err := r.JITFunc()
				Expect(url).To(Equal("resp"))
				Expect(err).To(BeNil())
			})

			rqstSvc.EXPECT().GetResponse(
				gomock.AssignableToTypeOf(&uuid.UUID{}),
			).DoAndReturn(func(uID *uuid.UUID) *requestM.RequestResponse {
				return &requestM.RequestResponse{
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

			r, e := target.BulkAsyncSessionRequest(
				[]func(*sessionM.Session) *requestM.Request{requestBuilder},
			)

			Expect(r).To(ConsistOf([][]byte{[]byte("stuff")}))
			Expect(e).To(ConsistOf(BeNil()))
		})
	})
})
