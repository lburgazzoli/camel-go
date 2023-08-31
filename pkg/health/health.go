package health

import (
	"context"
	"errors"
	"log"
	"net/http"
	"path"
	"sync"
	"sync/atomic"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type Check func() error

func New(address string, prefix string, logger *zap.Logger) *Service {
	s := Service{}

	s.router = gin.New()
	s.router.Use(ginzap.Ginzap(logger.Named("health"), time.RFC3339, true))
	s.router.GET(path.Join(prefix, "/ready"), s.ready)
	s.router.GET(path.Join(prefix, "/live"), s.live)

	s.srv = &http.Server{
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Addr:              address,
		Handler:           s.router,
	}

	return &s
}

type Service struct {
	lock    sync.Mutex
	running atomic.Bool
	router  *gin.Engine
	srv     *http.Server

	checksMutex     sync.RWMutex
	livenessChecks  map[string]Check
	readinessChecks map[string]Check
}

func (s *Service) Start(context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.running.CompareAndSwap(false, true) {
		go func() {
			err := s.srv.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("listen: %s\n", err)
			}
		}()
	}

	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.running.CompareAndSwap(true, false) {
		return s.srv.Shutdown(ctx)
	}

	return nil
}

func (s *Service) AddLivenessCheck(name string, check Check) {
	s.checksMutex.Lock()
	defer s.checksMutex.Unlock()

	s.livenessChecks[name] = check
}

func (s *Service) AddReadinessCheck(name string, check Check) {
	s.checksMutex.Lock()
	defer s.checksMutex.Unlock()

	s.readinessChecks[name] = check
}

func (s *Service) ready(c *gin.Context) {
	s.handle(c, s.readinessChecks)
}
func (s *Service) live(c *gin.Context) {
	s.handle(c, s.livenessChecks)
}

func (s *Service) collectChecks(checks map[string]Check, resultsOut map[string]string, statusOut *int) {
	s.checksMutex.RLock()
	defer s.checksMutex.RUnlock()

	for name, check := range checks {
		if err := check(); err != nil {
			*statusOut = http.StatusServiceUnavailable
			resultsOut[name] = err.Error()
		} else {
			resultsOut[name] = "OK"
		}
	}
}

func (s *Service) handle(c *gin.Context, checks ...map[string]Check) {
	checkResults := make(map[string]string)
	status := http.StatusOK

	for _, checks := range checks {
		s.collectChecks(checks, checkResults, &status)
	}

	switch c.Query("full") {
	case "true":
		c.JSON(status, gin.H{
			"status": "OK",
			"data":   checkResults,
		})
	default:
		c.JSON(status, gin.H{
			"status": "OK",
		})
	}
}
