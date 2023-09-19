// //go:build components_mqtt_v3 || components_all

package v3

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/lburgazzoli/camel-go/pkg/logger"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func init() {
	// TODO: do some investigation if there is a way to set per client
	//       logging instead of a global one
	mqtt.DEBUG = &mqttLogger{level: slog.LevelDebug}
	mqtt.WARN = &mqttLogger{level: slog.LevelWarn}
	mqtt.ERROR = &mqttLogger{level: slog.LevelError}
}

type mqttLogger struct {
	once   sync.Once
	level  slog.Level
	logger *slog.Logger
}

func (l *mqttLogger) Println(v ...interface{}) {
	l.once.Do(func() {
		l.logger = logger.WithGroup(Scheme)
	})

	switch l.level {
	case slog.LevelDebug:
		l.logger.Debug(fmt.Sprint(v...))
	case slog.LevelInfo:
		l.logger.Info(fmt.Sprint(v...))
	case slog.LevelWarn:
		l.logger.Warn(fmt.Sprint(v...))
	case slog.LevelError:
		l.logger.Error(fmt.Sprint(v...))
	}
}
func (l *mqttLogger) Printf(format string, v ...interface{}) {
	l.once.Do(func() {
		l.logger = logger.WithGroup(Scheme)
	})

	switch l.level {
	case slog.LevelDebug:
		l.logger.Debug(fmt.Sprintf(format, v...))
	case slog.LevelInfo:
		l.logger.Info(fmt.Sprintf(format, v...))
	case slog.LevelWarn:
		l.logger.Warn(fmt.Sprintf(format, v...))
	case slog.LevelError:
		l.logger.Error(fmt.Sprintf(format, v...))
	}
}
