////go:build components_mqtt_v3 || components_all

package v3

import (
	"fmt"
	"log/slog"
	"sync"

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
	m      sync.Mutex
	level  slog.Level
	logger *slog.Logger
}

func (l *mqttLogger) Println(v ...interface{}) {
	l.m.Lock()
	defer l.m.Unlock()

	if l.logger == nil {
		l.logger = slog.Default().WithGroup(Scheme)
	}

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
	l.m.Lock()
	defer l.m.Unlock()

	if l.logger == nil {
		l.logger = slog.Default().WithGroup(Scheme)
	}

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
