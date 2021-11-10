package session_service_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/JackStillwell/GoRez/internal/session_service"
	m "github.com/JackStillwell/GoRez/internal/session_service/models"
)

var _ = Describe("Service", func() {
	Context("NewSessionService", func() {
		It("should fail to create a session with more existing sessions than max sessions", func() {
			svc, err := session_service.NewSessionService(0, []*m.Session{{}})

			Expect(svc).To(BeNil())
			Expect(err.Error()).To(ContainSubstring("cannot create a session service"))
		})

		It("should add existing sessions to available sessions", func() {
			existingSession := &m.Session{
				Key: "123",
			}
			svc, err := session_service.NewSessionService(10, []*m.Session{existingSession})
			Expect(err).To(BeNil())

			sessChan := make(chan *m.Session, 10)
			svc.ReserveSession(1, sessChan)
			Eventually(sessChan).Should(Receive(Equal(existingSession)))
		})
	})

	Context("ReserveSession", func() {
		It("should return the requested number of sessions when avaialble", func() {
			svc, err := session_service.NewSessionService(5, []*m.Session{
				{Key: "123"},
				{Key: "123"},
				{Key: "123"},
				{Key: "123"},
				{Key: "123"},
			})

			Expect(err).To(BeNil())

			sessChan := make(chan *m.Session, 5)
			svc.ReserveSession(5, sessChan)

			for i := 0; i < 5; i++ {
				Eventually(sessChan).Should(Receive(Equal(&m.Session{Key: "123"})))
			}
		})
	})

	Context("ReleaseSession", func() {
		It("should make released sessions available", func() {
			svc, err := session_service.NewSessionService(1, []*m.Session{})
			Expect(err).To(BeNil())

			sessChan := make(chan *m.Session, 1)
			go svc.ReserveSession(1, sessChan)
			Consistently(sessChan).ShouldNot(Receive())

			sess := &m.Session{Key: "123"}
			svc.ReleaseSession([]*m.Session{sess})
			Eventually(sessChan).Should(Receive(Equal(sess)))
		})
	})

	Context("BadSession", func() {
		It("should remove one bad session from available sessions", func() {
			badSess := &m.Session{Key: "bad"}
			svc, err := session_service.NewSessionService(1, []*m.Session{badSess})
			Expect(err).To(BeNil())

			sessChan := make(chan *m.Session, 1)
			go svc.ReserveSession(1, sessChan)
			Eventually(sessChan).Should(Receive(Equal(badSess)))

			svc.BadSession([]*m.Session{badSess})
			go svc.ReserveSession(1, sessChan)
			Consistently(sessChan).ShouldNot(Receive())
		})

		It("should remove multiple bad sessions from available sessions", func() {
			badSess := &m.Session{Key: "bad"}
			goodSess1 := &m.Session{Key: "bad"}
			goodSess2 := &m.Session{Key: "bad"}
			svc, err := session_service.NewSessionService(3, []*m.Session{
				goodSess1,
				badSess,
				goodSess2,
			})
			Expect(err).To(BeNil())

			sessChan := make(chan *m.Session, 1)
			go svc.ReserveSession(3, sessChan)
			Eventually(sessChan).Should(Receive(Equal(goodSess1)))
			Eventually(sessChan).Should(Receive(Equal(badSess)))
			Eventually(sessChan).Should(Receive(Equal(goodSess2)))

			svc.ReleaseSession([]*m.Session{goodSess1})
			svc.BadSession([]*m.Session{badSess})
			svc.ReleaseSession([]*m.Session{goodSess2})

			go svc.ReserveSession(2, sessChan)
			Eventually(sessChan).Should(Receive(Equal(goodSess1)))
			Eventually(sessChan).Should(Receive(Equal(goodSess2)))
		})
	})
})