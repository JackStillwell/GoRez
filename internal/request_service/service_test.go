package request_service

import (
	"io"
	"net/http"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	m "github.com/JackStillwell/GoRez/internal/request_service/models"

	mock "github.com/JackStillwell/GoRez/internal/request_service/mocks"
)

var _ = Describe("Service", func() {
	var (
		ctrl *gomock.Controller
		rM   *requestManager
		rS   *requestService
	)

	uniqueId := uuid.New()

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Request", func() {
		var httpGet *mock.MockHTTPGet

		BeforeEach(func() {
			httpGet = mock.NewMockHTTPGet(ctrl)
			requester := NewTestRequester(httpGet)
			rM = NewTestRequestManager(5, requester)
			rS = &requestService{requester, rM}
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
				response = rS.Request(request)
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
				response = rS.Request(request)
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

				response = rS.Request(request)
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

				response = rS.Request(request)
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
			rM = NewTestRequestManager(5, requester)
			rS = &requestService{requester, rM}
		})

		Context("Is called with a request", func() {
			request := &m.Request{
				Id:      &uniqueId,
				JITArgs: []interface{}{"one", "two"},
				JITBuild: func([]interface{}) (string, error) {
					return "", nil
				},
			}

			It("should issue the request", func() {
				response := &m.RequestResponse{
					Id: &uniqueId,
				}
				requester.EXPECT().Request(request).Return(response).Times(1)
				rS.MakeRequest(request)
				Eventually(rM.responses, time.Second, time.Millisecond).
					Should(ContainElement(response))
			})
		})
	})

	Describe("GetResponse", func() {
		var requester *mock.MockRequester

		BeforeEach(func() {
			requester = mock.NewMockRequester(ctrl)
			rM = NewTestRequestManager(5, requester)
			rS = &requestService{requester, rM}
		})

		Context("Is called after a request has been made", func() {
			request := &m.Request{
				Id: &uniqueId,
			}

			requestResponse := &m.RequestResponse{
				Id: &uniqueId,
			}

			It("should return the response", func() {
				requester.EXPECT().Request(request).Return(requestResponse).Times(1)
				rS.MakeRequest(request)
				responseChan := make(chan *m.RequestResponse, 1)
				rS.GetResponse(&uniqueId, responseChan)
				response := <-responseChan
				Expect(response).To(Equal(requestResponse))
			})
		})

		Context("Can handle multiple listeners", func() {

			numRequests := 5
			IDs := make([]*uuid.UUID, numRequests)
			for idx := range IDs {
				uid := uuid.New()
				IDs[idx] = &uid
			}

			It("should return the response", func() {
				for _, ID := range IDs {
					requester.EXPECT().Request(
						&m.Request{
							Id: ID,
						},
					).Return(&m.RequestResponse{
						Id: ID,
					}).Times(1)
				}

				for _, ID := range IDs {
					rS.MakeRequest(&m.Request{Id: ID})
				}
				responseChan := make(chan *m.RequestResponse, numRequests)

				for _, ID := range IDs {
					rS.GetResponse(ID, responseChan)
				}

				responses := make([]*m.RequestResponse, numRequests)
				expecteds := make([]*m.RequestResponse, numRequests)
				for idx, ID := range IDs {
					expected := &m.RequestResponse{Id: ID}
					expecteds[idx] = expected
					response := <-responseChan
					responses[idx] = response
				}

				Expect(expecteds).To(ContainElements(responses))
			})
		})
	})
})
