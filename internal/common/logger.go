package common

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// InitLogger initializes the global logger
func InitLogger(level string) error {
	config := zap.NewProductionConfig()

	// Set log level
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Configure encoder
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.StacktraceKey = ""

	var err error
	logger, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// GetLogger returns the global logger
func GetLogger() *zap.Logger {
	if logger == nil {
		// Fallback to default logger
		logger = zap.NewExample()
	}
	return logger
}

// LoggerFromContext extracts logger from context or returns global logger
func LoggerFromContext(ctx context.Context) *zap.Logger {
	if ctxLogger, ok := ctx.Value("logger").(*zap.Logger); ok {
		return ctxLogger
	}
	return GetLogger()
}

// ContextWithLogger adds logger to context
func ContextWithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, "logger", logger)
}

// LogInfo logs info message
func LogInfo(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// LogError logs error message
func LogError(msg string, err error, fields ...zap.Field) {
	allFields := append(fields, zap.Error(err))
	GetLogger().Error(msg, allFields...)
}

// LogDebug logs debug message
func LogDebug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// LogWarn logs warning message
func LogWarn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Sync flushes any buffered log entries
func Sync() {
	if logger != nil {
		logger.Sync()
	}
}

// Close closes the logger
func Close() {
	if logger != nil {
		logger.Sync()
	}
}

// SetupLogger initializes logger based on environment
func SetupLogger() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	if err := InitLogger(logLevel); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}
