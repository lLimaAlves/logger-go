# logger-go

A structured, Prometheus-friendly logger for Go services — built on top of [go.uber.org/zap](https://github.com/uber-go/zap).

- **JSON output** in production (Loki/Prometheus-ready field names)
- **Colored console output** in development (auto-detected via `APP_ENV`)
- **Typed field constructors** — no stringly typed key/value pairs
- **Singleton + instance modes** — use the global singleton or create isolated loggers
- **Request-scoped child loggers** via `With()`
- Zero external dependencies beyond `go.uber.org/zap`

---

## Installation

```bash
go get github.com/lLimaAlves/logger-go
```

Requires Go 1.21+.

---

## Quick Start

```go
package main

import logger "github.com/lLimaAlves/logger-go"

func main() {
    // Initialise once at startup.
    logger.Init(logger.Config{
        Level:   logger.LevelInfo,
        Service: "my-service",
        Version: "1.0.0",
        Env:     "production",   // or leave empty to read APP_ENV
    })

    logger.Info("server started", logger.String("addr", ":8080"))
    logger.Warn("rate limit near threshold", logger.Int("requests", 980))
}
```

**JSON output (production)**
```json
{"level":"info","time":"2026-04-09T12:00:00.000Z","caller":"main/main.go:12","function":"main.main","msg":"server started","service":"my-service","version":"1.0.0","env":"production","addr":":8080"}
```

**Console output (development)**
```
2026-04-09T12:00:00.000Z  INFO  main/main.go:12  main.main  server started  {"service": "my-service", "env": "development", "addr": ":8080"}
```

---

## Configuration

```go
type Config struct {
    Level   Level  // LevelDebug | LevelInfo | LevelWarn | LevelError (default: LevelInfo)
    Service string // added to every log record as "service"
    Version string // added to every log record as "version" (optional)
    Env     string // added as "env"; reads APP_ENV if empty
}
```

The output format is selected automatically:

| `APP_ENV` / `Env` | Format |
|---|---|
| `production`, `prod`, `staging` | JSON |
| `development`, `dev`, or empty | Colored console |

---

## Usage

### Singleton (recommended for most services)

```go
// main.go — initialise once
logger.Init(logger.Config{Service: "api", Env: "production"})

// anywhere in the codebase
logger.Info("user created", logger.String("user_id", "u-123"))
logger.Error("db query failed", logger.Error(err), logger.String("query", "SELECT ..."))
```

### Isolated instance

```go
log := logger.New(logger.Config{
    Level:   logger.LevelDebug,
    Service: "worker",
})

log.Debug("job picked up", logger.Int("job_id", 7))
log.Info("job done", logger.Duration("elapsed", elapsed))
```

### Request-scoped child logger

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    log := logger.With(
        logger.String("request_id", r.Header.Get("X-Request-ID")),
        logger.String("user_id", userIDFromCtx(r.Context())),
    )

    log.Info("request received", logger.String("path", r.URL.Path))
    // All subsequent log calls in this scope carry request_id and user_id.
}
```

## Prometheus / Loki field names

Log records use lowercase snake_case keys by design, matching Prometheus label conventions:

```
level, time, caller, function, msg, service, version, env, stacktrace
```

This makes log-based alerting and metric extraction with LogQL/PromQL straightforward:

```logql
# Count error logs per service in the last 5 minutes
count_over_time({service="api"} | json | level="error" [5m])
```

---

## Log Levels

| Level | Use when |
|---|---|
| `LevelDebug` | Detailed diagnostic info, development only |
| `LevelInfo` | Normal operations, key lifecycle events |
| `LevelWarn` | Degraded or unexpected state, recoverable |
| `LevelError` | Failures requiring attention; stack trace attached automatically |


## License

MIT
