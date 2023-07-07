package gorez_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	gorez "github.com/JackStillwell/GoRez/pkg"
	c "github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"

	authMock "github.com/JackStillwell/GoRez/internal/auth/mocks"
	"github.com/JackStillwell/GoRez/internal/base"

	"github.com/JackStillwell/GoRez/internal/request"
	"github.com/JackStillwell/GoRez/internal/session"
)

var _ = Describe("ApiUtil", func() {
	Describe("Integrated Unit Tests", func() {
		var (
			ctrl       *gomock.Controller
			testServer *httptest.Server

			authSvc *authMock.MockService

			target i.APIUtil
			b      = base.NewService(zap.NewNop())
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			testServer = httptest.NewServer(nil)

			authSvc = authMock.NewMockService(ctrl)
			rqstSvc := request.NewService(3, b)
			sesnSvc := session.NewService(3, nil, b)

			hiRezC := c.NewHiRezConstants()
			hiRezC.SmiteURLBase = testServer.URL

			target = gorez.NewAPIUtil(hiRezC, authSvc, rqstSvc, sesnSvc)
		})

		AfterEach(func() {
			testServer.Close()
			ctrl.Finish()
		})

		Context("CreateSession", func() {
			// FIXME: don't know why this is failing, will need to fix
			It("should return requested sessions", func() {
				authSvc.EXPECT().GetID().Return("id").Times(3)
				authSvc.EXPECT().GetTimestamp(gomock.AssignableToTypeOf(time.Time{})).
					Return("timestamp").Times(3)
				authSvc.EXPECT().GetSignature(c.CreateSession, "timestamp").Return("signature").
					Times(3)

				testServer.Config.Handler = http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						defer GinkgoRecover()

						Expect(r.URL.Path).To(Equal("/createsessionjson/id/signature/timestamp"))
						w.WriteHeader(http.StatusOK)
						w.Write([]byte("{\"ret_msg\": \"Approved\", \"session_id\": \"session_id\"}"))
					},
				)

				retMsg := "Approved"
				sessionId := "session_id"
				sess := &m.Session{RetMsg: &retMsg, SessionID: &sessionId}
				sessions, errs := target.CreateSession(3)
				Expect(errs).To(ConsistOf(BeNil(), BeNil(), BeNil()))
				Expect(sessions).To(ConsistOf(sess, sess, sess))
			})
		})
	})
})
