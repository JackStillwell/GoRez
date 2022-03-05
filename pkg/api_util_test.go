package gorez_test

import (
	"net/http"
	"net/http/httptest"

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
			hiRezC.SmiteURLBase = testServer.URL + "/"

			target = gorez.NewAPIUtil(hiRezC, authSvc, rqstSvc, sesnSvc)
		})

		AfterEach(func() {
			testServer.Close()
			ctrl.Finish()
		})

		Context("CreateSession", func() {

		})
	})
})
