// File: logger/logger.go
// Package logger: tiny zap-based logger with fixed conventions.
// All comments are in English (per user's preference).
package logger

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log is the global sugared logger with fixed fields (service, env).
	Log *zap.SugaredLogger

	serviceName string
	envName     string
)

// Init initializes global logger with fixed conventions and common fields.
// level: "debug" | "info" | "warn" | "error" | ...
// service: stable service/binary name (e.g., "orders-api").
// env: environment tag (e.g., "prod" | "stage" | "dev").
func Init(level, service, env string) error {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		return fmt.Errorf("invalid log level %q: %w", level, err)
	}
	serviceName, envName = service, env

	// Encoder config: JSON, RFC3339 time, stable keys.
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = func(t time.Time, pa zapcore.PrimitiveArrayEncoder) {
		pa.AppendString(t.Format(time.RFC3339))
	}
	encCfg.TimeKey = "ts"     // timestamp
	encCfg.MessageKey = "msg" // message
	encCfg.CallerKey = "caller"

	encoder := zapcore.NewJSONEncoder(encCfg)
	ws := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(encoder, ws, zapLevel)

	// Add caller and stacktrace from error level.
	z := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	// Build sugared logger with fixed fields.
	l := z.Sugar().With(
		"service", serviceName,
		"env", envName,
	)

	// Expose globally and replace zap's global.
	Log = l
	zap.ReplaceGlobals(z)

	return nil
}

// New returns a child logger with the required "module" field.
// If Init wasn't called, it creates a sensible default and proceeds.
func New(module string) *zap.SugaredLogger {
	if Log == nil {
		_ = Init("info", "unknown-service", "dev")
	}
	return Log.With("module", module)
}

// Sync flushes any buffered log entries (call on graceful shutdown).
func Sync() error {
	if Log == nil {
		return nil
	}
	return Log.Sync()
}
