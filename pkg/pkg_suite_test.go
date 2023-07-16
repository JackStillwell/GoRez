package gorez_test

import (
	"strings"
	"testing"

	rqstM "github.com/JackStillwell/GoRez/internal/request/models"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPkg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pkg Suite")
}

func NewRequstURLContainsMatcher(expected string) gomock.Matcher {
	return gomock.GotFormatterAdapter(RequestURLContainsGotFormatter{}, &RequestURLContainsMatcher{
		Expected: expected,
	})
}

type RequestURLContainsMatcher struct {
	Expected string
	errMsg   string
}

func (m *RequestURLContainsMatcher) Matches(x any) bool {
	switch x := x.(type) {
	case *rqstM.Request:
		url, err := x.JITFunc()
		if err != nil {
			m.errMsg = err.Error()
			return false
		}
		return strings.Contains(url, m.Expected)
	default:
		m.errMsg = "unexpected type"
		return false
	}
}

func (m RequestURLContainsMatcher) String() string {
	if m.errMsg != "" {
		return "error creating URL: " + m.errMsg
	}
	return "expected request to contain " + m.Expected
}

type RequestURLContainsGotFormatter struct{}

func (m RequestURLContainsGotFormatter) Got(x any) string {
	switch x := x.(type) {
	case *rqstM.Request:
		url, err := x.JITFunc()
		if err != nil {
			return "error creating URL: " + err.Error()
		}
		return url
	default:
		return "unexpected type"
	}
}
