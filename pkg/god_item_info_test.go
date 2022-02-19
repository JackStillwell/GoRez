package gorez_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

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
	requestM "github.com/JackStillwell/GoRez/internal/request_service/models"

	session "github.com/JackStillwell/GoRez/internal/session_service"
	sessionI "github.com/JackStillwell/GoRez/internal/session_service/interfaces"
	sessionMock "github.com/JackStillwell/GoRez/internal/session_service/mocks"
	sessionM "github.com/JackStillwell/GoRez/internal/session_service/models"
)

var _ = Describe("GodItemInfo", func() {
	var (
		hiRezConsts c.HiRezConstants

		target i.GodItemInfo
	)

	Describe("IntegratedUnitTest", func() {
		var (
			authSvc authI.AuthService
			rqstSvc requestI.RequestService
			sesnSvc sessionI.SessionService
		)

		BeforeEach(func() {
			authSvc = auth.NewAuthService(authM.Auth{
				ID:  "id",
				Key: "key",
			})
			rqstSvc = request.NewRequestService(1)
			sesnSvc = session.NewSessionService(1, []*sessionM.Session{{}})
		})

		Context("singleRequest via GetGods", func() {
			It("should return an error from requesting a response", func() {
				testServer := httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusInternalServerError)
					}))
				defer testServer.Close()

				hiRezConsts.SmiteURLBase = testServer.URL + "/"

				target = gorez.NewGodItemInfo(hiRezConsts, rqstSvc, authSvc, sesnSvc)

				_, err := target.GetGods()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("requesting response"))
				Expect(err.Error()).To(ContainSubstring(fmt.Sprint(
					http.StatusInternalServerError,
				)))
			})

			It("should return an error from unmarshaling a response body", func() {
				testServer := httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusOK)
					}))
				defer testServer.Close()

				hiRezConsts.SmiteURLBase = testServer.URL + "/"

				target = gorez.NewGodItemInfo(hiRezConsts, rqstSvc, authSvc, sesnSvc)

				_, err := target.GetGods()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("marshaling response"))
			})

			It("should return success if no errors occur", func() {
				testServer := httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusOK)
						rw.Write([]byte("[{},{}]"))
					}))
				defer testServer.Close()

				hiRezConsts.SmiteURLBase = testServer.URL + "/"

				target = gorez.NewGodItemInfo(hiRezConsts, rqstSvc, authSvc, sesnSvc)

				gods, err := target.GetGods()
				Expect(err).ToNot(HaveOccurred())
				Expect(gods).To(HaveLen(2))
			})
		})

		Context("multirequest via GetGodRecItems", func() {
			It("should remove a bad session", func() {
				testServer := httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusInternalServerError)
						rw.Write([]byte("session"))
					}))
				defer testServer.Close()

				hiRezConsts.SmiteURLBase = testServer.URL + "/"

				target = gorez.NewGodItemInfo(hiRezConsts, rqstSvc, authSvc, sesnSvc)

				_, errs := target.GetGodRecItems([]int{0})
				Expect(errs).To(HaveLen(1))
				Expect(errs[0].Error()).To(And(
					ContainSubstring("request"),
					ContainSubstring(fmt.Sprint(http.StatusInternalServerError)),
					ContainSubstring("session"),
				))

				Expect(sesnSvc.GetAvailableSessions()).To(HaveLen(0))
			})

			It("should return a good session", func() {
				testServer := httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusInternalServerError)
					}))
				defer testServer.Close()

				hiRezConsts.SmiteURLBase = testServer.URL + "/"

				target = gorez.NewGodItemInfo(hiRezConsts, rqstSvc, authSvc, sesnSvc)

				_, errs := target.GetGodRecItems([]int{0})
				Expect(errs).To(HaveLen(1))
				Expect(errs[0].Error()).To(And(
					ContainSubstring("request"),
					ContainSubstring(fmt.Sprint(http.StatusInternalServerError)),
				))

				Expect(sesnSvc.GetAvailableSessions()).To(HaveLen(1))
			})

			It("should return an error from requesting a response", func() {
				testServer := httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusInternalServerError)
					}))
				defer testServer.Close()

				hiRezConsts.SmiteURLBase = testServer.URL + "/"

				target = gorez.NewGodItemInfo(hiRezConsts, rqstSvc, authSvc, sesnSvc)

				_, errs := target.GetGodRecItems([]int{0})
				Expect(errs).To(HaveLen(1))
				Expect(errs[0].Error()).To(And(
					ContainSubstring("request"),
					ContainSubstring(fmt.Sprint(http.StatusInternalServerError)),
				))
			})

			It("should return an error from marshaling a response", func() {
				testServer := httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusOK)
						rw.Write([]byte(""))
					}))
				defer testServer.Close()

				hiRezConsts.SmiteURLBase = testServer.URL + "/"

				target = gorez.NewGodItemInfo(hiRezConsts, rqstSvc, authSvc, sesnSvc)

				_, errs := target.GetGodRecItems([]int{0})
				Expect(errs).To(HaveLen(1))
				Expect(errs[0].Error()).To(And(
					ContainSubstring("marshaling response"),
				))
			})
		})
	})

	Describe("UnitTests", func() {
		var (
			ctrl *gomock.Controller

			authSvc *authMock.MockAuthService
			rqstSvc *requestMock.MockRequestService
			sesnSvc *sessionMock.MockSessionService
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			authSvc = authMock.NewMockAuthService(ctrl)
			rqstSvc = requestMock.NewMockRequestService(ctrl)
			sesnSvc = sessionMock.NewMockSessionService(ctrl)

			hiRezConsts = c.NewHiRezConstants()

			target = gorez.NewGodItemInfo(hiRezConsts, rqstSvc, authSvc, sesnSvc)
		})

		Context("singleRequest via GetItems", func() {
			It("should pass through an error with the request", func() {
				s := &sessionM.Session{Key: "key"}
				var sesnChan chan *sessionM.Session
				sesnSvc.EXPECT().ReserveSession(1, gomock.AssignableToTypeOf(sesnChan)).Do(
					func(_ int, c chan *sessionM.Session) {
						sesnChan = c
						sesnChan <- s
					},
				)
				sesnSvc.EXPECT().ReleaseSession([]*sessionM.Session{s})

				authSvc.EXPECT().GetID().Return("id")

				rqstSvc.EXPECT().Request(gomock.AssignableToTypeOf(&requestM.Request{})).Return(
					&requestM.RequestResponse{
						Err: errors.New("boom"),
					})

				_, err := target.GetItems()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(And(
					ContainSubstring("requesting response"),
					ContainSubstring("boom"),
				))
			})

			It("should pass through an error unmarshaling the response", func() {
				s := &sessionM.Session{Key: "key"}
				var sesnChan chan *sessionM.Session
				sesnSvc.EXPECT().ReserveSession(1, gomock.AssignableToTypeOf(sesnChan)).Do(
					func(_ int, c chan *sessionM.Session) {
						sesnChan = c
						sesnChan <- s
					},
				)
				sesnSvc.EXPECT().ReleaseSession([]*sessionM.Session{s})

				authSvc.EXPECT().GetID().Return("id")

				rqstSvc.EXPECT().Request(gomock.AssignableToTypeOf(&requestM.Request{})).Return(
					&requestM.RequestResponse{
						Err:  nil,
						Resp: []byte(""),
					})

				_, err := target.GetItems()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(And(
					ContainSubstring("marshaling response"),
				))
			})
		})
	})
})
