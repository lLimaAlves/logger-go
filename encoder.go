package logger

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// jsonEncoderConfig returns the zapcore encoder config used for JSON output.
// Field names are chosen to be compatible with Prometheus/Loki label conventions
// (lowercase, snake_case).
func jsonEncoderConfig(loc *time.Location) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "function",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // prometheus-friendly: debug, info, warn, error
		EncodeTime:     iso8601InLocation(loc),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// consoleEncoderConfig returns the zapcore encoder config used for human-readable
// colored output in development.
func consoleEncoderConfig(loc *time.Location) zapcore.EncoderConfig {
	cfg := jsonEncoderConfig(loc)
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return cfg
}

// iso8601InLocation returns a TimeEncoder that formats timestamps as ISO 8601
// in the given location. Falls back to time.UTC when loc is nil.
func iso8601InLocation(loc *time.Location) zapcore.TimeEncoder {
	if loc == nil {
		loc = time.UTC
	}
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.In(loc).Format("2006-01-02T15:04:05.000Z07:00"))
	}
}
