package request_service_test

import (
	"io"
	"net/http"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	s "github.com/JackStillwell/GoRez/internal/request_service"
	i "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	m "github.com/JackStillwell/GoRez/internal/request_service/models"

	mock "github.com/JackStillwell/GoRez/internal/request_service/mocks"
)

var _ = Describe("Service", func() {
	var (
		ctrl           *gomock.Controller
		requestService i.RequestService
		httpGet        *mock.MockHTTPGet
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		httpGet = mock.NewMockHTTPGet(ctrl)
		requestService = s.NewTestRequestService(5, httpGet)
	})

	AfterEach(func() {
		ctrl.Finish()
		requestService.Close()
	})

	Describe("Request", func() {
		Context("Encounters a URL build error", func() {
			uniqueId := uuid.New()
			request := &m.Request{
				Id: &uniqueId,
				JITBuild: func([]interface{}) (string, error) {
					return "", errors.New("unexpected")
				},
			}

			It("should return the error", func() {
				response := requestService.Request(request)
				Expect(response.Err).ToNot(BeNil())
				Expect(response.Err.Error()).To(Equal("building requesturl: unexpected"))
			})

			It("should have the same Id as the request", func() {
				response := requestService.Request(request)
				Expect(response.Id).To(Equal(request.Id))
			})

			It("should have a nil Resp", func() {
				response := requestService.Request(request)
				Expect(response.Resp).To(BeNil())
			})
		})

		Context("Encounters an HTTPGet error", func() {
			uniqueId := uuid.New()
			request := &m.Request{
				Id: &uniqueId,
				JITBuild: func([]interface{}) (string, error) {
					return "", nil
				},
			}

			BeforeEach(func() {
				httpGet.EXPECT().Get("").Return(nil, errors.New("unexpected"))
			})

			It("should return the error", func() {
				response := requestService.Request(request)
				Expect(response.Err).ToNot(BeNil())
				Expect(response.Err.Error()).To(Equal("getting response: unexpected"))
			})

			It("should have the same Id as the request", func() {
				response := requestService.Request(request)
				Expect(response.Id).To(Equal(request.Id))
			})

			It("should have a nil Resp", func() {
				response := requestService.Request(request)
				Expect(response.Resp).To(BeNil())
			})
		})

		Context("Encounters a body reading error", func() {
			uniqueId := uuid.New()
			request := &m.Request{
				Id: &uniqueId,
				JITBuild: func([]interface{}) (string, error) {
					return "", nil
				},
			}

			BeforeEach(func() {
				mRC := mock.NewMockReadCloser(ctrl)
				mRC.EXPECT().Read(gomock.Any()).Return(0, errors.New("unexpected"))
				mRC.EXPECT().Close()
				httpGet.EXPECT().Get("").Return(&http.Response{
					Body: mRC,
				}, nil)
			})

			It("should return the error", func() {
				response := requestService.Request(request)
				Expect(response.Err).ToNot(BeNil())
				Expect(response.Err.Error()).To(Equal("reading body: unexpected"))
			})

			It("should have the same Id as the request", func() {
				response := requestService.Request(request)
				Expect(response.Id).To(Equal(request.Id))
			})

			It("should have a nil Resp", func() {
				response := requestService.Request(request)
				Expect(response.Resp).To(BeNil())
			})
		})

		Context("Encounters no errors", func() {
			uniqueId := uuid.New()
			request := &m.Request{
				Id: &uniqueId,
				JITBuild: func([]interface{}) (string, error) {
					return "", nil
				},
			}

			BeforeEach(func() {
				mRC := mock.NewMockReadCloser(ctrl)
				mRC.EXPECT().Read(gomock.Any()).
					SetArg(0, []byte("hello world")).
					Return(11, io.EOF)
				mRC.EXPECT().Close()
				httpGet.EXPECT().Get("").Return(&http.Response{
					Body: mRC,
				}, nil)
			})

			It("should have a nil error", func() {
				response := requestService.Request(request)
				Expect(response.Err).To(BeNil())
			})

			It("should have the same Id as the request", func() {
				response := requestService.Request(request)
				Expect(response.Id).To(Equal(request.Id))
			})

			It("should have the Resp returned by httpget", func() {
				response := requestService.Request(request)
				Expect(response.Resp).To(Equal([]byte("hello world")))
			})
		})
	})
})
