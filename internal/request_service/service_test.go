package request_service_test

import (
	"github.com/golang/mock/gomock"
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

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		requestService = s.NewMockRequestService(5, mock.NewMockHTTPGet(ctrl))
	})

	Describe("Request", func() {
		Context("Encounters a URL build error", func() {
			var request *m.Request

			BeforeEach(func() {
				request = &m.Request{
					JITBuild: func(urlBuilder string, buildArgs []interface{}) (string, error) {
						return "", errors.New("unexpected")
					},
				}
			})

			It("should return the error", func() {
				response := requestService.Request(request)
				Expect(response.Err).ToNot(BeNil())
				Expect(response.Err.Error()).To(Equal("building requesturl: unexpected"))
			})
		})
	})

	AfterEach(func() {
		ctrl.Finish()
	})
})
