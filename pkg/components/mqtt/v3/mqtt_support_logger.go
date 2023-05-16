////go:build components_mqtt_v3 || components_all

package v3

import (
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lburgazzoli/camel-go/pkg/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	// TODO: do some investigation if there is a way to set per client
	//       logging instead of a global one
	mqtt.DEBUG = &mqttLogger{level: zapcore.DebugLevel}
	mqtt.WARN = &mqttLogger{level: zapcore.WarnLevel}
	mqtt.ERROR = &mqttLogger{level: zapcore.ErrorLevel}
}

type mqttLogger struct {
	m      sync.Mutex
	level  zapcore.Level
	logger *zap.SugaredLogger
}

func (l *mqttLogger) Println(v ...interface{}) {
	l.m.Lock()
	defer l.m.Unlock()

	if core.L == nil {
		return
	}
	if l.logger == nil {
		l.logger = core.L.Named(Scheme).Sugar()
	}

	switch l.level {
	case zapcore.DebugLevel:
		l.logger.Debug(v...)
	case zapcore.InfoLevel:
		l.logger.Info(v...)
	case zapcore.WarnLevel:
		l.logger.Warn(v...)
	case zapcore.ErrorLevel:
		l.logger.Error(v...)
	}
}
func (l *mqttLogger) Printf(format string, v ...interface{}) {
	l.m.Lock()
	defer l.m.Unlock()

	if core.L == nil {
		return
	}
	if l.logger == nil {
		l.logger = core.L.Named(Scheme).Sugar()
	}

	switch l.level {
	case zapcore.DebugLevel:
		l.logger.Debugf(format, v...)
	case zapcore.InfoLevel:
		l.logger.Infof(format, v...)
	case zapcore.WarnLevel:
		l.logger.Warnf(format, v...)
	case zapcore.ErrorLevel:
		l.logger.Errorf(format, v...)
	}
}
