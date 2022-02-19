package internal_test

import (
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	rSMocks "github.com/JackStillwell/GoRez/internal/request_service/mocks"
	rSM "github.com/JackStillwell/GoRez/internal/request_service/models"
	sSMocks "github.com/JackStillwell/GoRez/internal/session_service/mocks"
	sSM "github.com/JackStillwell/GoRez/internal/session_service/models"

	"github.com/JackStillwell/GoRez/internal"
)

var _ = Describe("Utils", func() {
	Context("BulkAsyncSessionRequest", func() {
		var (
			ctrl *gomock.Controller

			rqstSvc *rSMocks.MockRequestService
			sesnSvc *sSMocks.MockSessionService
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())

			rqstSvc = rSMocks.NewMockRequestService(ctrl)
			sesnSvc = sSMocks.NewMockSessionService(ctrl)
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should return each response", func() {
			sess := &sSM.Session{
				Key: "fake",
			}
			sessChan := make(chan *sSM.Session, 1)
			sesnSvc.EXPECT().ReserveSession(1, gomock.AssignableToTypeOf(sessChan)).
				Do(func(_ int, sessChan chan *sSM.Session) {
					sessChan <- sess
				})

			requestBuilder := func(s *sSM.Session) *rSM.Request {
				return &rSM.Request{
					JITArgs: []interface{}{s},
					JITBuild: func([]interface{}) (string, error) {
						return "resp", nil
					},
				}
			}

			request := requestBuilder(sess)
			rqstSvc.EXPECT().MakeRequest(gomock.AssignableToTypeOf(request)).Do(func(r *rSM.Request) {
				request.Id = r.Id
				Expect(r.JITArgs).To(Equal(request.JITArgs))
				Expect(r.JITBuild).To(BeAssignableToTypeOf(request.JITBuild))
			})

			respChan := make(chan *rSM.RequestResponse, 1)
			rqstSvc.EXPECT().GetResponse(
				gomock.AssignableToTypeOf(&uuid.UUID{}),
				gomock.AssignableToTypeOf(respChan),
			).Do(func(uID *uuid.UUID, respChan chan *rSM.RequestResponse) {
				respChan <- &rSM.RequestResponse{
					Id:   uID,
					Err:  nil,
					Resp: []byte("stuff"),
				}
			})

			sesnSvc.EXPECT().ReleaseSession(gomock.AssignableToTypeOf([]*sSM.Session{})).Do(
				func(sessions []*sSM.Session) {
					Expect(sessions).To(ConsistOf(sess))
				},
			)

			responsesChan := make(chan [][]byte, 1)
			errorsChan := make(chan []error, 1)
			go func() {
				defer GinkgoRecover()

				r, e := internal.BulkAsyncSessionRequest(
					rqstSvc,
					sesnSvc,
					[]func(*sSM.Session) *rSM.Request{requestBuilder},
				)
				responsesChan <- r
				errorsChan <- e
			}()

			Eventually(responsesChan).Should(Receive(ConsistOf([][]byte{[]byte("stuff")})))
			Eventually(errorsChan).Should(Receive(&[]error{nil}))
		})
	})
})
