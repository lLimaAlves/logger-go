package logger

import "go.uber.org/zap/zapcore"

// jsonEncoderConfig returns the zapcore encoder config used for JSON output.
// Field names are chosen to be compatible with Prometheus/Loki label conventions
// (lowercase, snake_case).
func jsonEncoderConfig() zapcore.EncoderConfig {
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
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// consoleEncoderConfig returns the zapcore encoder config used for human-readable
// colored output in development.
func consoleEncoderConfig() zapcore.EncoderConfig {
	cfg := jsonEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return cfg
}
