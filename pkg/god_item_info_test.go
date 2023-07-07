package gorez_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	gorez "github.com/JackStillwell/GoRez/pkg"
	c "github.com/JackStillwell/GoRez/pkg/constants"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	mock "github.com/JackStillwell/GoRez/pkg/mocks"

	"github.com/JackStillwell/GoRez/internal/auth"
	authI "github.com/JackStillwell/GoRez/internal/auth/interfaces"
	authM "github.com/JackStillwell/GoRez/internal/auth/models"
	"github.com/JackStillwell/GoRez/internal/base"

	"github.com/JackStillwell/GoRez/internal/request"
	requestI "github.com/JackStillwell/GoRez/internal/request/interfaces"

	"github.com/JackStillwell/GoRez/internal/session"
	sessionI "github.com/JackStillwell/GoRez/internal/session/interfaces"
	sessionM "github.com/JackStillwell/GoRez/internal/session/models"
)

var _ = Describe("GodItemInfo", func() {
	var (
		hiRezConsts c.HiRezConstants

		target i.GodItemInfo
	)

	Describe("IntegratedUnitTest", func() {
		var (
			util    i.GorezUtil
			authSvc authI.Service
			rqstSvc requestI.Service
			sesnSvc sessionI.Service

			b = base.NewService(zap.NewNop())
		)

		BeforeEach(func() {
			authSvc = auth.NewService(authM.Auth{
				ID:  "id",
				Key: "key",
			}, b)
			rqstSvc = request.NewService(1, b)
			sesnSvc = session.NewService(1, []*sessionM.Session{{}}, b)

			util = gorez.NewGorezUtil(authSvc, rqstSvc, sesnSvc)
		})

		Context("singleRequest via GetGods", func() {
			It("should return an error from requesting a response", func() {
				testServer := httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusInternalServerError)
					}))
				defer testServer.Close()

				hiRezConsts.SmiteURLBase = testServer.URL + "/"

				target = gorez.NewGodItemInfo(hiRezConsts, util)

				_, err := target.GetGods()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("requesting response"))
				Expect(err.Error()).To(ContainSubstring(fmt.Sprint(
					http.StatusInternalServerError,
				)))
			})

			It("should return an error from single response ret_msg", func() {
				testServer := httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusOK)
						rw.Write([]byte("{\"ret_msg\": \"issue\"}"))
					}))
				defer testServer.Close()

				hiRezConsts.SmiteURLBase = testServer.URL + "/"

				target = gorez.NewGodItemInfo(hiRezConsts, util)

				_, err := target.GetGods()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ret_msg: issue"))
			})

			It("should return an error from array response ret_msg", func() {
				testServer := httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusOK)
						rw.Write([]byte("["))
						rw.Write([]byte("{\"ret_msg\": \"\"},"))
						rw.Write([]byte("{\"ret_msg\": \"issue\"}"))
						rw.Write([]byte("]"))
					}))
				defer testServer.Close()

				hiRezConsts.SmiteURLBase = testServer.URL + "/"

				target = gorez.NewGodItemInfo(hiRezConsts, util)

				_, err := target.GetGods()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ret_msg 1: issue"))
			})

			It("should return success if no errors occur", func() {
				testServer := httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusOK)
						rw.Write([]byte("[{},{}]"))
					}))
				defer testServer.Close()

				hiRezConsts.SmiteURLBase = testServer.URL + "/"

				target = gorez.NewGodItemInfo(hiRezConsts, util)

				gods, err := target.GetGods()
				Expect(err).ToNot(HaveOccurred())
				Expect(gods).To(Equal([]byte("[{},{}]")))
			})
		})

		Context("multirequest via GetGodRecItems", func() {
			It("should remove a bad session", func(ctx SpecContext) {
				testServer := httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusInternalServerError)
						rw.Write([]byte("session"))
					}))
				defer testServer.Close()

				hiRezConsts.SmiteURLBase = testServer.URL + "/"

				target = gorez.NewGodItemInfo(hiRezConsts, util)

				_, errs := target.GetGodRecItems([]int{0})
				Expect(errs).To(HaveLen(1))
				Expect(errs[0].Error()).To(And(
					ContainSubstring("request"),
					ContainSubstring(fmt.Sprint(http.StatusInternalServerError)),
					ContainSubstring("session"),
				))

				Expect(sesnSvc.GetAvailableSessions()).To(HaveLen(0))
			}, NodeTimeout(5*time.Second))

			It("should return a good session", func() {
				testServer := httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusInternalServerError)
					}))
				defer testServer.Close()

				hiRezConsts.SmiteURLBase = testServer.URL + "/"

				target = gorez.NewGodItemInfo(hiRezConsts, util)

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

				target = gorez.NewGodItemInfo(hiRezConsts, util)

				_, errs := target.GetGodRecItems([]int{0})
				Expect(errs).To(HaveLen(1))
				Expect(errs[0].Error()).To(And(
					ContainSubstring("request"),
					ContainSubstring(fmt.Sprint(http.StatusInternalServerError)),
				))
			})
		})
	})

	Describe("UnitTests", func() {
		var (
			ctrl *gomock.Controller

			util *mock.MockGorezUtil
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			util = mock.NewMockGorezUtil(ctrl)

			hiRezConsts = c.HiRezConstants{}

			target = gorez.NewGodItemInfo(hiRezConsts, util)
		})

		Context("singleRequest via GetItems", func() {
			It("should pass through an error with the request", func() {
				util.EXPECT().SingleRequest(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("boom"))

				_, err := target.GetItems()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(And(
					ContainSubstring("boom"),
				))
			})
		})
	})
})
