package request_service_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/JackStillwell/GoRez/internal/request_service"
	i "github.com/JackStillwell/GoRez/internal/request_service/interfaces"
	m "github.com/JackStillwell/GoRez/internal/request_service/models"
)

var _ = Describe("Service", func() {
	var (
		uniqueId uuid.UUID

		target i.RequestService
	)

	BeforeEach(func() {
		uniqueId = uuid.New()

		target = request_service.NewRequestService(1)
	})

	Describe("Request", func() {
		Context("Encounters a URL build error", func() {
			request := &m.Request{
				Id: &uniqueId,
				JITBuild: func([]interface{}) (string, error) {
					return "", errors.New("unexpected")
				},
			}

			var response *m.RequestResponse
			BeforeEach(func() {
				response = target.Request(request)
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
			var (
				request  *m.Request
				response *m.RequestResponse
			)

			BeforeEach(func() {
				request = &m.Request{
					Id: &uniqueId,
					JITBuild: func([]interface{}) (string, error) {
						return "invalidurl", nil
					},
				}
				response = target.Request(request)
			})

			It("should return the error", func() {
				Expect(response.Err).To(HaveOccurred())
				Expect(response.Err.Error()).To(
					ContainSubstring("getting response"),
				)
			})

			It("should have the same Id as the request", func() {
				Expect(response.Id).To(Equal(request.Id))
			})

			It("should have a nil Resp", func() {
				Expect(response.Resp).To(BeNil())
			})
		})

		Context("Encounters no errors", func() {
			var (
				server   *httptest.Server
				request  *m.Request
				response *m.RequestResponse
			)

			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusOK)
						rw.Write([]byte("hello world"))
					}))
				request = &m.Request{
					Id: &uniqueId,
					JITBuild: func([]interface{}) (string, error) {
						return server.URL, nil
					},
				}
				response = target.Request(request)
			})

			AfterEach(func() {
				server.Close()
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
	/*
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
	*/
})
