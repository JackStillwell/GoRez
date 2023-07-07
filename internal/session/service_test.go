package session_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/JackStillwell/GoRez/internal/base"
	"github.com/JackStillwell/GoRez/internal/session"
	i "github.com/JackStillwell/GoRez/internal/session/interfaces"
	m "github.com/JackStillwell/GoRez/internal/session/models"
)

var _ = Describe("Service", func() {
	var (
		b = base.NewService(zap.NewNop())
	)
	Context("NewSessionService", func() {
		It("should fail to create a session with more existing sessions than max sessions", func() {
			Expect(func() {
				session.NewService(0, []*m.Session{{}}, b)
			}).To(PanicWith(
				"cannot create a session service with capacity 0 and 1 existing sessions",
			))
		})

		It("should add existing sessions to available sessions", func() {
			existingSession := &m.Session{
				Key: "123",
			}
			var svc i.Service
			Expect(func() {
				svc = session.NewService(10, []*m.Session{existingSession}, b)
			}).ToNot(Panic())

			sessChan := make(chan *m.Session, 10)
			svc.ReserveSession(1, sessChan)
			Eventually(sessChan).Should(Receive(Equal(existingSession)))
		})
	})

	Context("GetAvailableSessions", func() {
		It("should return all currently available sessions", func() {
			numSessions := 5
			svc := session.NewService(numSessions, []*m.Session{
				{Key: "123"},
				{Key: "123"},
				{Key: "123"},
				{Key: "123"},
				{Key: "123"},
			}, b)

			Expect(svc.GetAvailableSessions()).To(HaveLen(numSessions))
		})
	})

	Context("ReserveSession", func() {
		It("should return the requested number of sessions when avaialble", func() {
			svc := session.NewService(5, []*m.Session{
				{Key: "123"},
				{Key: "123"},
				{Key: "123"},
				{Key: "123"},
				{Key: "123"},
			}, b)

			sessChan := make(chan *m.Session, 5)
			svc.ReserveSession(5, sessChan)

			for i := 0; i < 5; i++ {
				Eventually(sessChan).Should(Receive(Equal(&m.Session{Key: "123"})))
			}
		})
	})

	Context("ReleaseSession", func() {
		It("should make released sessions available", func() {
			svc := session.NewService(1, []*m.Session{}, b)

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
			svc := session.NewService(1, []*m.Session{badSess}, b)

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
			svc := session.NewService(3, []*m.Session{
				goodSess1,
				badSess,
				goodSess2,
			}, b)

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
