// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

// Logger --
type Logger struct {
	Level  string `yaml:"level"`
	Writer string `yaml:"writer"`
}

// LogConfiguration --
type LogConfiguration struct {
	Loggers map[string]Logger
}

// ==========================
//
// Global
//
// ==========================

// Configuration --
var Configuration LogConfiguration

// RootLogger --
var RootLogger zerolog.Logger

// Log --
func Log(level zerolog.Level, format string, args ...interface{}) {
	RootLogger.WithLevel(level).Msgf(format, args)
}

// ==========================
//
//
//
// ==========================

// New --
func New(name string) zerolog.Logger {
	if cfg, ok := Configuration.Loggers[name]; ok {
		level := toLevel(cfg.Level)
		writer := toWriter(cfg.Writer)

		if level == zerolog.Disabled {
			return zerolog.Nop()
		}

		l := zerolog.New(writer).With().Timestamp().Str("logger", name).Logger()
		l.Hook(discardHook{threshold: level})

		return l
	}

	return zerolog.New(os.Stdout).With().Timestamp().Str("logger", name).Logger()
}

// ==========================
//
// Helpers
//
// ==========================

func toLevel(level string) zerolog.Level {
	if strings.EqualFold("debug", level) {
		return zerolog.DebugLevel
	}
	if strings.EqualFold("info", level) {
		return zerolog.InfoLevel
	}
	if strings.EqualFold("warn", level) {
		return zerolog.WarnLevel
	}
	if strings.EqualFold("error", level) {
		return zerolog.ErrorLevel
	}
	if strings.EqualFold("fatal", level) {
		return zerolog.FatalLevel
	}
	if strings.EqualFold("panic", level) {
		return zerolog.PanicLevel
	}
	if strings.EqualFold("disabled", level) {
		return zerolog.Disabled
	}

	return zerolog.NoLevel
}

func toWriter(writer string) io.Writer {
	if strings.EqualFold("stdout", writer) {
		return os.Stdout
	}
	if strings.EqualFold("stderr", writer) {
		return os.Stderr
	}
	if strings.EqualFold("disabled", writer) {
		return ioutil.Discard
	}

	// TODO: files
	return os.Stdout
}

type discardHook struct {
	threshold zerolog.Level
}

func (h discardHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level < h.threshold || h.threshold >= zerolog.NoLevel {
		e.Discard()
	}
}

// ==========================
//
// Initialization
//
// ==========================

func init() {
	RootLogger = New("root")

	// default configuration
	Configuration = LogConfiguration{
		Loggers: make(map[string]Logger),
	}
}
