package logger

import (
	"time"

	"go.uber.org/zap"
)


// Use the constructor functions below (String, Int, Bool, …) rather than constructing zap.Field directly. 
type Field = zap.Field

// --- String-like ---------------------------------------------------------

// String adds a string field.
func String(key, value string) Field { return zap.String(key, value) }

// Stringer adds a field whose value is the result of calling .String() on v.
func Stringer(key string, v interface{ String() string }) Field {
	return zap.Stringer(key, v)
}

// Error logs err under the key "error". If err is nil, the field is a no-op.
// Prefer returning errors over logging them mid-stack; log once at the boundary.
func Error(err error) Field { return zap.Error(err) }

// NamedError logs err under a custom key.
func NamedError(key string, err error) Field { return zap.NamedError(key, err) }

// --- Numeric -------------------------------------------------------------

// Int adds an int field.
func Int(key string, value int) Field { return zap.Int(key, value) }

// Int64 adds an int64 field.
func Int64(key string, value int64) Field { return zap.Int64(key, value) }

// Uint adds a uint field.
func Uint(key string, value uint) Field { return zap.Uint(key, value) }

// Float64 adds a float64 field.
func Float64(key string, value float64) Field { return zap.Float64(key, value) }

// --- Boolean -------------------------------------------------------------

// Bool adds a bool field.
func Bool(key string, value bool) Field { return zap.Bool(key, value) }

// --- Duration / Time -----------------------------------------------------

// Duration adds a time.Duration field encoded as fractional seconds.
func Duration(key string, value time.Duration) Field { return zap.Duration(key, value) }

// Time adds a time.Time field encoded in ISO 8601.
func Time(key string, value time.Time) Field { return zap.Time(key, value) }

// --- Any (escape hatch) --------------------------------------------------

// Any adds a field with an arbitrary value. Use a typed constructor when
// possible; this variant is provided for external types and dynamic values.
func Any(key string, value any) Field { return zap.Any(key, value) }
