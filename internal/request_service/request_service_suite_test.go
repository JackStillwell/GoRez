package request_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRequestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RequestService Suite")
}
