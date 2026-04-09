package logger_test

import (
	"errors"
	"time"

	logger "github.com/lLimaAlves/logger-go"
)

// ExampleInit demonstrates initialising the singleton at application startup.
func ExampleInit() {
	logger.Init(logger.Config{
		Level:   logger.LevelInfo,
		Service: "my-service",
		Version: "1.2.3",
		Env:     "production",
	})
}

// ExampleNew demonstrates creating an isolated logger (no singleton).
func ExampleNew() {
	log := logger.New(logger.Config{
		Level:   logger.LevelDebug,
		Service: "worker",
		Env:     "development",
	})

	log.Debug("processing job", logger.Int("job_id", 42))
	log.Info("job complete", logger.Duration("elapsed", 120*time.Millisecond))
}

// ExampleLogger_With demonstrates building request-scoped child loggers.
func ExampleLogger_With() {
	log := logger.New(logger.Config{Service: "api", Env: "production"})

	// Attach request-scoped fields once; all subsequent calls carry them.
	reqLog := log.With(
		logger.String("request_id", "abc-123"),
		logger.String("user_id", "u-456"),
	)

	reqLog.Info("handling request",
		logger.String("method", "POST"),
		logger.String("path", "/orders"),
	)
	reqLog.Error("payment failed", logger.Error(errors.New("card declined")))
}

// ExampleGetLogger demonstrates the package-level singleton helpers.
func ExampleGetLogger() {
	// Initialise once in main.
	logger.Init(logger.Config{Service: "demo", Env: "development"})

	// Package-level helpers are available anywhere without passing the logger.
	logger.Info("application boot complete")
	logger.Warn("config not found, using defaults",
		logger.String("config_path", "/etc/app.yaml"),
	)
	logger.Debug("debug details", logger.Any("payload", map[string]any{"key": "value"}))
}
