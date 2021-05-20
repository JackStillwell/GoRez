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
	)

	uniqueId := uuid.New()

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
		requestService.Close()
	})

	Describe("Request", func() {
		var httpGet *mock.MockHTTPGet

		BeforeEach(func() {
			httpGet = mock.NewMockHTTPGet(ctrl)
			requester := s.NewTestRequester(httpGet)
			requestService = s.NewTestRequestService(5, requester)
		})

		Context("Encounters a URL build error", func() {
			request := &m.Request{
				Id: &uniqueId,
				JITBuild: func([]interface{}) (string, error) {
					return "", errors.New("unexpected")
				},
			}

			var response *m.RequestResponse
			BeforeEach(func() {
				response = requestService.Request(request)
			})

			It("should return the error", func() {
				Expect(response.Err).ToNot(BeNil())
				Expect(response.Err.Error()).To(Equal("building requesturl: unexpected"))
			})

			It("should have the same Id as the request", func() {
				Expect(response.Id).To(Equal(request.Id))
			})

			It("should have a nil Resp", func() {
				Expect(response.Resp).To(BeNil())
			})
		})

		Context("Encounters an HTTPGet error", func() {
			request := &m.Request{
				Id: &uniqueId,
				JITBuild: func([]interface{}) (string, error) {
					return "", nil
				},
			}

			var response *m.RequestResponse
			BeforeEach(func() {
				httpGet.EXPECT().Get("").Return(nil, errors.New("unexpected"))
				response = requestService.Request(request)
			})

			It("should return the error", func() {
				Expect(response.Err).ToNot(BeNil())
				Expect(response.Err.Error()).To(Equal("getting response: unexpected"))
			})

			It("should have the same Id as the request", func() {
				Expect(response.Id).To(Equal(request.Id))
			})

			It("should have a nil Resp", func() {
				Expect(response.Resp).To(BeNil())
			})
		})

		Context("Encounters a body reading error", func() {
			request := &m.Request{
				Id: &uniqueId,
				JITBuild: func([]interface{}) (string, error) {
					return "", nil
				},
			}

			var response *m.RequestResponse
			BeforeEach(func() {
				mRC := mock.NewMockReadCloser(ctrl)
				mRC.EXPECT().Read(gomock.Any()).Return(0, errors.New("unexpected"))
				mRC.EXPECT().Close()
				httpGet.EXPECT().Get("").Return(&http.Response{
					Body: mRC,
				}, nil)

				response = requestService.Request(request)
			})

			It("should return the error", func() {
				Expect(response.Err).ToNot(BeNil())
				Expect(response.Err.Error()).To(Equal("reading body: unexpected"))
			})

			It("should have the same Id as the request", func() {
				Expect(response.Id).To(Equal(request.Id))
			})

			It("should have a nil Resp", func() {
				Expect(response.Resp).To(BeNil())
			})
		})

		Context("Encounters no errors", func() {
			request := &m.Request{
				Id: &uniqueId,
				JITBuild: func([]interface{}) (string, error) {
					return "", nil
				},
			}

			var response *m.RequestResponse
			BeforeEach(func() {
				mRC := mock.NewMockReadCloser(ctrl)
				mRC.EXPECT().Read(gomock.Any()).
					SetArg(0, []byte("hello world")).
					Return(11, io.EOF)
				mRC.EXPECT().Close()
				httpGet.EXPECT().Get("").Return(&http.Response{
					Body: mRC,
				}, nil)

				response = requestService.Request(request)
			})

			It("should have a nil error", func() {
				Expect(response.Err).To(BeNil())
			})

			It("should have the same Id as the request", func() {
				Expect(response.Id).To(Equal(request.Id))
			})

			It("should have the Resp returned by httpget", func() {
				Expect(response.Resp).To(Equal([]byte("hello world")))
			})
		})
	})

	Describe("MakeRequest", func() {
		var requester *mock.MockRequester

		BeforeEach(func() {
			requester = mock.NewMockRequester(ctrl)
			requestService = s.NewTestRequestService(5, requester)
		})

		Context("Is called with a request", func() {
			/*
				request := &m.Request{
					Id:      &uniqueId,
					JITArgs: []interface{}{"one", "two"},
					JITBuild: func([]interface{}) (string, error) {
						return "", nil
					},
				}
			*/

			It("should issue the request", func() {
				// Have no way to make the test wait for the call
				/*
					requester.EXPECT().Request(request).Times(1)
					requestService.MakeRequest(request)
				*/
			})
		})
	})

	Describe("GetRequest", func() {
		var requester *mock.MockRequester

		BeforeEach(func() {
			requester = mock.NewMockRequester(ctrl)
			requestService = s.NewTestRequestService(5, requester)
		})

		Context("Is called after a request has been made", func() {
			// Need to figure out async testing
			/*
				request := &m.Request{
					Id: &uniqueId,
				}

				requestResponse := &m.RequestResponse{
					Id: &uniqueId,
				}

				It("should return the response", func() {
					requester.EXPECT().Request(request).Return(requestResponse).Times(1)
					requestService.MakeRequest(request)
					Expect(requestService.GetResponse(&uniqueId)).To(Equal(requestResponse))
				})
			*/
		})
	})
})
