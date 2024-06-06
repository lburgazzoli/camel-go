package dapr

import (
	"log/slog"
	"sync"

	"github.com/dapr/go-sdk/service/common"

	daprd "github.com/dapr/go-sdk/service/http"
)

var s = daprdSvc{
	cnt: 0,
	svc: daprd.NewService(Address()),
	log: slog.Default().With(slog.String("subsystem", "daprd")),
}

func Start() error {
	return s.Start()
}

func Stop() error {
	return s.Stop()
}

func AddTopicEventHandler(sub *common.Subscription, fn common.TopicEventHandler) error {
	return s.AddTopicEventHandler(sub, fn)
}

type daprdSvc struct {
	mu  sync.Mutex
	cnt uint32
	svc common.Service
	log *slog.Logger
}

func (s *daprdSvc) Start() error {
	if s == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: better ref counter
	if s.cnt == 0 {
		s.log.Info("staring")

		if err := s.svc.Start(); err != nil {
			return err
		}

		s.cnt++
	}

	return nil
}

func (s *daprdSvc) Stop() error {
	if s == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// TODO: better ref counter
	if s.cnt == 1 {
		s.log.Info("stopping")

		if err := s.svc.Stop(); err != nil {
			return err
		}

		s.cnt--
	}

	return nil
}

func (s *daprdSvc) AddTopicEventHandler(sub *common.Subscription, fn common.TopicEventHandler) error {
	return s.svc.AddTopicEventHandler(sub, fn)
}
