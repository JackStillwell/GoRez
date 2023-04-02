package request_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/JackStillwell/GoRez/internal/request"
	i "github.com/JackStillwell/GoRez/internal/request/interfaces"
	m "github.com/JackStillwell/GoRez/internal/request/models"
)

var _ = Describe("Service", func() {
	var (
		uniqueId uuid.UUID

		target i.Service
	)

	BeforeEach(func() {
		uniqueId = uuid.New()

		target = request.NewService(1)
	})

	Describe("Request", func() {
		Context("Encounters a URL build error", func() {
			request := &m.Request{
				Id: &uniqueId,
				JITFunc: func() (string, error) {
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
					JITFunc: func() (string, error) {
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

		Context("Encounters a response code error", func() {
			var (
				server   *httptest.Server
				request  *m.Request
				response *m.RequestResponse
			)

			BeforeEach(func() {
				server = httptest.NewServer(http.HandlerFunc(
					func(rw http.ResponseWriter, r *http.Request) {
						rw.WriteHeader(http.StatusInternalServerError)
						rw.Write([]byte("boom"))
					}))
				request = &m.Request{
					Id: &uniqueId,
					JITFunc: func() (string, error) {
						return server.URL, nil
					},
				}
				response = target.Request(request)
			})

			AfterEach(func() {
				server.Close()
			})

			It("should return an error including status code and body", func() {
				Expect(response.Err).To(HaveOccurred())
				Expect(response.Err.Error()).To(
					ContainSubstring("status code"),
					ContainSubstring(fmt.Sprint(http.StatusInternalServerError)),
					ContainSubstring("boom"),
				)
			})

			It("should have the same Id as the request", func() {
				Expect(response.Id).To(Equal(request.Id))
			})

			It("should have a Resp containing the returned body", func() {
				Expect(response.Resp).To(Equal([]byte("boom")))
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
					JITFunc: func() (string, error) {
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

	Describe("MakeRequest", func() {
		var (
			server  *httptest.Server
			request *m.Request
		)

		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(
				func(rw http.ResponseWriter, r *http.Request) {
					rw.WriteHeader(http.StatusOK)
					rw.Write([]byte("hello world"))
				}))
			request = &m.Request{
				Id: &uniqueId,
				JITFunc: func() (string, error) {
					return server.URL, nil
				},
			}
		})

		AfterEach(func() {
			server.Close()
		})

		Context("Is called with a request", func() {
			It("should issue the request", func() {
				target.MakeRequest(request)
				responseChan := make(chan *m.RequestResponse, 1)
				target.GetResponse(&uniqueId, responseChan)
				Eventually(responseChan).Should(Receive())
			})
		})
	})

	Describe("GetResponse", func() {
		var (
			server *httptest.Server
		)

		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(
				func(rw http.ResponseWriter, r *http.Request) {
					rw.WriteHeader(http.StatusOK)
				}))
		})

		AfterEach(func() {
			server.Close()
		})

		Context("Can handle multiple listeners", func() {

			numRequests := 5
			IDs := make([]*uuid.UUID, numRequests)
			for idx := range IDs {
				uid := uuid.New()
				IDs[idx] = &uid
			}

			BeforeEach(func() {
				target = request.NewService(numRequests)
			})

			It("should return the response", func() {
				for _, ID := range IDs {
					target.MakeRequest(&m.Request{
						Id: ID,
						JITFunc: func() (string, error) {
							return server.URL, nil
						},
					})
				}
				responseChan := make(chan *m.RequestResponse, numRequests)

				for _, ID := range IDs {
					target.GetResponse(ID, responseChan)
				}

				responses := make([]*m.RequestResponse, numRequests)
				expecteds := make([]*m.RequestResponse, numRequests)
				for idx, ID := range IDs {
					expected := &m.RequestResponse{Id: ID, Resp: []byte{}}
					expecteds[idx] = expected
					response := <-responseChan
					responses[idx] = response
				}

				Expect(expecteds).To(ContainElements(responses))
			})
		})
	})
})
