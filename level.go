package logger

import "go.uber.org/zap/zapcore"

type Level string

const (
	// LevelDebug is the most verbose level; for development diagnostics.
	LevelDebug Level = "debug"
	// LevelInfo is the default level for normal operational events.
	LevelInfo Level = "info"
	// LevelWarn indicates degraded or unexpected states that are recoverable.
	LevelWarn Level = "warn"
	// LevelError indicates failures that require attention.
	LevelError Level = "error"
)

// toZap converts a Level to its zapcore equivalent.
// Unrecognised values fall back to InfoLevel.
func (l Level) toZap() zapcore.Level {
	switch l {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
