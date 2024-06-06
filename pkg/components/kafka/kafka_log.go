package kafka

import (
	"context"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"
)

type klog struct {
	delegate *slog.Logger
}

func (l *klog) Level() kgo.LogLevel {
	switch {
	case l.delegate.Enabled(context.TODO(), slog.LevelDebug):
		return kgo.LogLevelDebug
	case l.delegate.Enabled(context.TODO(), slog.LevelInfo):
		return kgo.LogLevelInfo
	case l.delegate.Enabled(context.TODO(), slog.LevelWarn):
		return kgo.LogLevelWarn
	case l.delegate.Enabled(context.TODO(), slog.LevelError):
		return kgo.LogLevelError
	default:
		return kgo.LogLevelNone
	}
}

func (l *klog) Log(level kgo.LogLevel, msg string, keyvals ...any) {
	switch level {
	case kgo.LogLevelDebug:
		l.delegate.Debug(msg, keyvals...)
	case kgo.LogLevelError:
		l.delegate.Error(msg, keyvals...)
	case kgo.LogLevelInfo:
		l.delegate.Info(msg, keyvals...)
	case kgo.LogLevelWarn:
		l.delegate.Warn(msg, keyvals...)
	default: // do nothing
	}
}
