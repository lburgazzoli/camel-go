package logger

import (
	"os"

	"github.com/rs/zerolog"
)

// New --
func New(logger string) zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Str("logger", logger).Logger()
}

// ==========================
//
// Global
//
// ==========================

var rootLogger zerolog.Logger

// Log --
func Log(level zerolog.Level, format string, args ...interface{}) {
	rootLogger.WithLevel(level).Msgf(format, args)
}

// ==========================
//
// Initialization
//
// ==========================

func init() {
	rootLogger = New("root")
}
