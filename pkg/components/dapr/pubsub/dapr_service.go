package pubsub

import (
	"sync"

	"github.com/dapr/go-sdk/service/common"

	daprd "github.com/dapr/go-sdk/service/http"
)

// TODO: better ref counter

func NewService(address string) *Service {
	return &Service{
		cnt: 0,
		svc: daprd.NewService(address),
	}
}

type Service struct {
	mu  sync.Mutex
	cnt uint32
	svc common.Service
}

func (s *Service) Start() error {
	if s == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.svc == nil && s.cnt == 0 {
		if err := s.svc.Start(); err != nil {
			return err
		}

		s.cnt++

	}

	return nil
}

func (s *Service) Stop() error {
	if s == nil {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.svc != nil && s.cnt == 1 {
		if err := s.svc.Stop(); err != nil {
			return err
		}

		s.cnt--
	}

	return nil
}

func (s *Service) AddTopicEventHandler(sub *common.Subscription, fn common.TopicEventHandler) error {
	return s.svc.AddTopicEventHandler(sub, fn)
}
