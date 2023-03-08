////go:build components_mqtt || components_all

package mqtt

import (
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
	level  zapcore.Level
	logger *zap.SugaredLogger
}

func (l *mqttLogger) Println(v ...interface{}) {
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
