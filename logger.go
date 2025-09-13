// File: logger/logger.go
// Package logger: tiny zap-based logger with fixed conventions.
// All comments are in English (per user's preference).
package logger

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
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

const (
	defaultLevel   = "info"
	defaultService = "unknown-service"
	defaultEnv     = "dev"
)

// Init initializes global logger with fixed conventions and common fields.
// level: "debug" | "info" | "warn" | "error" | ...
// service: stable service/binary name (e.g., "orders-api").
// env: environment tag (e.g., "prod" | "stage" | "dev").
//
// If level is empty, it will fallback to LOG_LEVEL env var, then to "info".
func Init(level, service, env string) error {
	// Resolve level: explicit arg > LOG_LEVEL > default
	lvlStr := strings.TrimSpace(level)
	if lvlStr == "" {
		lvlStr = strings.TrimSpace(os.Getenv("LOG_LEVEL"))
	}
	if lvlStr == "" {
		lvlStr = defaultLevel
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(strings.ToLower(lvlStr))); err != nil {
		return fmt.Errorf("invalid log level %q: %w", lvlStr, err)
	}

	// Resolve service/env with safe defaults
	if strings.TrimSpace(service) == "" {
		service = defaultService
	}
	if strings.TrimSpace(env) == "" {
		env = defaultEnv
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
	encCfg.LevelKey = "level" // make sure it stays "level"

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

// ensureDefaultInit makes sure there is a usable global logger.
// It is used when New() is called before Init().
func ensureDefaultInit() {
	if Log != nil {
		return
	}
	// Safe lazy default: level=info, service/env defaults.
	_ = Init(defaultLevel, defaultService, defaultEnv)
}

// New returns a child logger with the required "module" field.
// If Init wasn't called, it creates a sensible default and proceeds.
func New(module string) *zap.SugaredLogger {
	ensureDefaultInit()
	return Log.With("module", module)
}

// Sync flushes any buffered log entries (call on graceful shutdown).
// EPIPE on Linux during stdout close is benign and will be ignored.
func Sync() error {
	if Log == nil {
		return nil
	}
	if err := Log.Sync(); err != nil {
		// Ignore EPIPE: happens when stdout/stderr is already closed (e.g., piped).
		if errors.Is(err, syscall.EPIPE) {
			return nil
		}
		return err
	}
	return nil
}
