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
