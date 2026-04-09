package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Use New to construct one or Init/GetLogger for the package-level singleton.
type Logger struct {
	zap   *zap.Logger
	level zap.AtomicLevel
}

// Logger Options.
type Config struct {
	// Level sets the minimum log severity. Defaults to LevelInfo.
	Level Level
	// Service is added to every log record as the "service" field.
	Service string
	// Version is added to every log record as the "version" field (optional).
	Version string
	// Env is added to every log record as the "env" field (optional).
	// When empty it reads APP_ENV from the environment.
	Env string
}

var (
	global   *Logger
	initOnce sync.Once
)

// Init initialises the package-level singleton.
// Call this once at startup before any logging occurs.
func Init(cfg Config) *Logger {
	initOnce.Do(func() {
		global = New(cfg)
	})
	return global
}

// GetLogger returns the package-level singleton, initialising it with
// sensible defaults if Init has not been called yet.
// Call Init before any logging to use a custom config; a GetLogger call
// that races Init will prevent Init from taking effect.
func GetLogger() *Logger {
	initOnce.Do(func() {
		global = New(Config{Level: LevelInfo})
	})
	return global
}

// New builds a Logger from cfg. It never panics; fatal build errors call os.Exit(1).
func New(cfg Config) *Logger {
	env := cfg.Env
	if env == "" {
		env = os.Getenv("APP_ENV")
	}

	isDev := env == "development" || env == "dev" || env == ""

	var enc zapcore.Encoder
	if isDev {
		enc = zapcore.NewConsoleEncoder(consoleEncoderConfig())
	} else {
		enc = zapcore.NewJSONEncoder(jsonEncoderConfig())
	}

	atomicLevel := zap.NewAtomicLevelAt(cfg.Level.toZap())
	sink := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(enc, sink, atomicLevel)

	base := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(1), // skip Logger method wrapper
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	// Attach default fields
	fields := make([]zap.Field, 0, 3)
	if cfg.Service != "" {
		fields = append(fields, zap.String("service", cfg.Service))
	}
	if cfg.Version != "" {
		fields = append(fields, zap.String("version", cfg.Version))
	}
	if env != "" {
		fields = append(fields, zap.String("env", env))
	}

	if len(fields) > 0 {
		base = base.With(fields...)
	}

	return &Logger{zap: base, level: atomicLevel}
}

// With returns a child Logger with the given fields pre-attached.
// Use this to create request-scoped or component-scoped loggers.
func (l *Logger) With(fields ...Field) *Logger {
	return &Logger{zap: l.zap.With(fields...), level: l.level}
}

// UpdateLevel atomically changes the minimum log level at runtime.
// Useful for toggling debug output without restarting the service.
func (l *Logger) UpdateLevel(level Level) {
	l.level.SetLevel(level.toZap())
}

// Zap returns the underlying *zap.Logger for interop with libraries that
// accept a *zap.Logger directly (e.g. grpc-zap, gorm logger adapters).
func (l *Logger) Zap() *zap.Logger {
	return l.zap
}

// --------------------------------------------------------------------------
// Logging methods
// --------------------------------------------------------------------------

// Debug logs a message at the debug level with optional structured fields.
func (l *Logger) Debug(msg string, fields ...Field) {
	l.zap.Debug(msg, fields...)
}

// Info logs a message at the info level with optional structured fields.
func (l *Logger) Info(msg string, fields ...Field) {
	l.zap.Info(msg, fields...)
}

// Warn logs a message at the warn level with optional structured fields.
func (l *Logger) Warn(msg string, fields ...Field) {
	l.zap.Warn(msg, fields...)
}

// Error logs a message at the error level with optional structured fields.
// A stack trace is automatically attached.
func (l *Logger) Error(msg string, fields ...Field) {
	l.zap.Error(msg, fields...)
}

// Fatal logs a message at the fatal level, then calls os.Exit(1).
// Use only in main or top-level bootstrapping code.
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.zap.Fatal(msg, fields...)
}

// --------------------------------------------------------------------------
// Package-level helpers (delegate to the singleton)
// --------------------------------------------------------------------------

// Debug logs at debug level using the singleton logger.
func Debug(msg string, fields ...Field) { GetLogger().Debug(msg, fields...) }

// Info logs at info level using the singleton logger.
func Info(msg string, fields ...Field) { GetLogger().Info(msg, fields...) }

// Warn logs at warn level using the singleton logger.
func Warn(msg string, fields ...Field) { GetLogger().Warn(msg, fields...) }

// ErrorLog logs at error level using the singleton logger.
// Named ErrorLog to avoid collision with the Error(err) Field constructor.
func ErrorLog(msg string, fields ...Field) { GetLogger().Error(msg, fields...) }

// Fatal logs at fatal level using the singleton logger, then exits.
func Fatal(msg string, fields ...Field) { GetLogger().Fatal(msg, fields...) }

// With returns a child of the singleton with pre-attached fields.
func With(fields ...Field) *Logger { return GetLogger().With(fields...) }
