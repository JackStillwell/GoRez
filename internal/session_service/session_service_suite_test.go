package session_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSessionService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SessionService Suite")
}
