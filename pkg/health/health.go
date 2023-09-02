package health

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"path"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPrefix   = ""
	DefaultAddress  = ":8081"
	DefaultPort     = 8081
	DefaultPortName = "health"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type Check func() error

func New(address string, prefix string, logger *slog.Logger) *Service {
	s := Service{}
	s.l = logger.WithGroup("health")

	s.router = gin.New()
	s.router.Use(s.log)
	s.router.GET(path.Join(prefix, "/health", "/ready"), s.ready)
	s.router.GET(path.Join(prefix, "/health", "/live"), s.live)

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
	l       *slog.Logger
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

func (s *Service) log(c *gin.Context) {
	start := time.Now()

	// some evil middlewares modify this values
	urlPath := c.Request.URL.Path
	urlQuery := c.Request.URL.RawQuery

	c.Next()

	end := time.Now()
	latency := end.Sub(start)

	fields := []any{
		slog.Int("status", c.Writer.Status()),
		slog.String("method", c.Request.Method),
		slog.String("path", urlPath),
		slog.String("query", urlQuery),
		slog.String("ip", c.ClientIP()),
		slog.String("user-agent", c.Request.UserAgent()),
		slog.Duration("latency", latency),
	}

	if trace.SpanFromContext(c.Request.Context()).SpanContext().HasTraceID() {
		fields = append(fields, trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String())
	}

	if len(c.Errors) > 0 {
		for _, e := range c.Errors.Errors() {
			s.l.Error(e, fields...)
		}
	} else {
		s.l.Info(urlPath, fields...)
	}
}
