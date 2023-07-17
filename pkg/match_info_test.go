package gorez_test

import (
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gorez "github.com/JackStillwell/GoRez/pkg"
	i "github.com/JackStillwell/GoRez/pkg/interfaces"
	m "github.com/JackStillwell/GoRez/pkg/models"

	auth "github.com/JackStillwell/GoRez/internal/auth"
	authI "github.com/JackStillwell/GoRez/internal/auth/interfaces"
	authM "github.com/JackStillwell/GoRez/internal/auth/models"

	rqstMocks "github.com/JackStillwell/GoRez/internal/request/mocks"
	rqstM "github.com/JackStillwell/GoRez/internal/request/models"

	sesnMocks "github.com/JackStillwell/GoRez/internal/session/mocks"
	sesnM "github.com/JackStillwell/GoRez/internal/session/models"
)

var _ = Describe("MatchInfo", func() {
	var (
		ctrl *gomock.Controller

		sesnSvc *sesnMocks.MockService
		authSvc authI.Service
		rqstSvc *rqstMocks.MockService

		matchInfo i.MatchInfo
	)

	BeforeEach(func(ctx SpecContext) {
		ctrl = gomock.NewController(GinkgoT())

		sesnSvc = sesnMocks.NewMockService(ctrl)
		authSvc = auth.NewService(authM.Auth{
			ID:  "1",
			Key: "123",
		}, nil)
		rqstSvc = rqstMocks.NewMockService(ctrl)

		matchInfo = gorez.NewMatchInfo(rqstSvc, authSvc, sesnSvc)
	}, NodeTimeout(time.Second))

	Context("GetMatchIdsByQueue", func() {
		FIt("should request each provided dateString queue combination", func(ctx SpecContext) {
			sesnSvc.EXPECT().ReserveSession(1, gomock.Any()).Do(
				func(_ int, c chan *sesnM.Session) {
					c <- &sesnM.Session{Key: "123"}
				}).Times(2)

			rqstSvc.EXPECT().MakeRequest(NewRequstURLContainsMatcher("20230716/0")).Times(1)
			rqstSvc.EXPECT().MakeRequest(NewRequstURLContainsMatcher("20230716/1")).Times(1)
			rqstSvc.EXPECT().GetResponse(gomock.Any()).DoAndReturn(
				func(uID *uuid.UUID) *rqstM.RequestResponse {
					return &rqstM.RequestResponse{
						Id:   uID,
						Err:  nil,
						Resp: []byte(`[{"Match": "123"}]`),
					}
				},
			).Times(2)
			sesnSvc.EXPECT().ReleaseSession(gomock.Any()).Times(2)

			results, errs := matchInfo.GetMatchIDsByQueue([]string{"20230716/0", "20230716/1"}, []m.QueueID{m.RankedConquest})
			Expect(errs).To(ConsistOf(BeNil(), BeNil()))

			expectedMatchID := "123"
			Expect(results).To(ConsistOf(&[]m.MatchIDWithQueue{
				{MatchID: m.MatchID{Match: &expectedMatchID}, QueueID: 451},
			}, &[]m.MatchIDWithQueue{
				{MatchID: m.MatchID{Match: &expectedMatchID}, QueueID: 451},
			}))
		}, SpecTimeout(time.Second*1))
	})
})
