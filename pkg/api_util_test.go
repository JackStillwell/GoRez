package gorez_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gorez "github.com/JackStillwell/GoRez/pkg"
	c "github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"

	authMock "github.com/JackStillwell/GoRez/internal/auth_service/mocks"

	request "github.com/JackStillwell/GoRez/internal/request_service"

	session "github.com/JackStillwell/GoRez/internal/session_service"
)

var _ = Describe("ApiUtil", func() {
	Describe("Integrated Unit Tests", func() {
		var (
			ctrl       *gomock.Controller
			testServer *httptest.Server

			serverFunc func(w http.ResponseWriter, r *http.Request)

			authSvc *authMock.MockAuthService

			target i.APIUtil
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			testServer := httptest.NewServer(http.HandlerFunc(serverFunc))

			authSvc = authMock.NewMockAuthService(ctrl)
			rqstSvc := request.NewRequestService(3)
			sesnSvc := session.NewSessionService(3, nil)

			hiRezC := c.NewHiRezConstants()
			hiRezC.SmiteURLBase = testServer.URL

			target = gorez.NewAPIUtil(hiRezC, authSvc, rqstSvc, sesnSvc)
		})

		AfterEach(func() {
			testServer.Close()
			ctrl.Finish()
		})

		Context("CreateSession", func() {
			FIt("should return requested sessions", func() {
				authSvc.EXPECT().GetID().Return("id").Times(3)
				authSvc.EXPECT().GetTimestamp(gomock.AssignableToTypeOf(time.Time{})).
					Return("timestamp").Times(3)
				authSvc.EXPECT().GetSignature(c.CreateSession, "timestamp").Return("signature").
					Times(3)

				serverFunc = func(w http.ResponseWriter, r *http.Request) {
					defer GinkgoRecover()

					Expect(r.URL.Path).To(Equal(""))
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("body"))
				}

				done := make(chan bool)
				go func(done chan bool) {
					defer GinkgoRecover()

					sessions, errs := target.CreateSession(3)
					Expect(errs).To(HaveLen(0))
					Expect(sessions).To(HaveLen(3))
					Expect(sessions).To(ConsistOf("", "", ""))
					done <- true
				}(done)

				select {
				case <-time.After(time.Second):
					Fail("timeout")
				case <-done:
					// nothing means the test passes
				}
			})
		})
	})
})
