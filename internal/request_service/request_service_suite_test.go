package request_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRequestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RequestService Suite")
}
