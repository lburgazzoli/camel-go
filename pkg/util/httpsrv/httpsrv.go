package httpsrv

import (
	"net/http"
	"time"
)

const (
	DefaultReadTimeout       = 1 * time.Second
	DefaultWriteTimeout      = 1 * time.Second
	DefaultIdleTimeout       = 30 * time.Second
	DefaultReadHeaderTimeout = 2 * time.Second
)

func New(address string, handler http.Handler) *http.Server {
	return &http.Server{
		ReadTimeout:       DefaultReadTimeout,
		WriteTimeout:      DefaultWriteTimeout,
		IdleTimeout:       DefaultIdleTimeout,
		ReadHeaderTimeout: DefaultReadHeaderTimeout,
		Addr:              address,
		Handler:           handler,
	}
}
