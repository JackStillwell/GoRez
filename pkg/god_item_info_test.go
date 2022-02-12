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

	auth "github.com/JackStillwell/GoRez/internal/auth_service"
	authI "github.com/JackStillwell/GoRez/internal/auth_service/interfaces"
	authMock "github.com/JackStillwell/GoRez/internal/auth_service/mocks"
	authM "github.com/JackStillwell/GoRez/internal/auth_service/models"

	request "github.com/JackStillwell/GoRez/internal/request_service"
	requestI "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	requestMock "github.com/JackStillwell/GoRez/internal/request_service/mocks"

	session "github.com/JackStillwell/GoRez/internal/session_service"
	sessionI "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
	sessionMock "github.com/JackStillwell/GoRez/internal/session_service/mocks"
	sessionM "github.com/JackStillwell/GoRez/internal/session_service/models"
)

var _ = Describe("GodItemInfo", func() {
	var (
		ctrl *gomock.Controller

		authSvc authI.AuthService
		rqstSvc requestI.RequestService
		sesnSvc sessionI.SessionService

		hiRezConsts c.HiRezConstants

		target i.GodItemInfo
	)

	Describe("IntegratedUnitTest", func() {
		BeforeEach(func() {
			authSvc = auth.NewAuthService(authM.Auth{})
			rqstSvc = request.NewRequestService(1)
			sesnSvc = session.NewSessionService(1, []*sessionM.Session{{}})

			hiRezConsts = c.NewHiRezConstants()

			target = gorez.NewGodItemInfo(hiRezConsts, rqstSvc, authSvc, sesnSvc)
		})

		FContext("GetGods", func() {
			It("should return an error from requesting a response", func() {
				testServer := httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusInternalServerError)
					}))
				hiRezConsts.SmiteURLBase = testServer.URL + "/"

				_, err := target.GetGods()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("requesting response"))
			})
		})
	})

	Describe("UnitTests", func() {
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			authSvc = authMock.NewMockAuthService(ctrl)
			rqstSvc = requestMock.NewMockRequestService(ctrl)
			sesnSvc = sessionMock.NewMockSessionService(ctrl)

			hiRezConsts = c.NewHiRezConstants()

			target = gorez.NewGodItemInfo(hiRezConsts, rqstSvc, authSvc, sesnSvc)
		})
	})
})
