package internal_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/JackStillwell/GoRez/internal"
)

var _ = Describe("Utils", func() {
	Context("Get", func() {
		It("should return an error from get", func() {
			testServer := httptest.NewServer(http.HandlerFunc(
				func(rw http.ResponseWriter, r *http.Request) {
					rw.WriteHeader(http.StatusInternalServerError)
				}))
			defer testServer.Close()

			_, err := internal.DefaultGetter{}.Get("invalidurl")
			Expect(err).To(HaveOccurred())
		})
	})
})
