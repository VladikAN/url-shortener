package logger

import (
	"github.com/vladikan/url-shortener/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger

// Init will setup logger
func Init(st *config.LogSettings) {
	cfg := zap.NewProductionConfig()
	cfg.Level = zap.NewAtomicLevelAt(toZapLevel(st.Level))

	lg, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	logger = lg
	sugar = logger.Sugar()
}

// Fatal logs message with Fatal level
func Fatal(msg string) {
	sugar.Fatal(msg)
}

// Info logs message with Info level
func Info(msg string) {
	sugar.Info(msg)
}

// Warn logs message with Warn level
func Warn(msg string) {
	sugar.Warn(msg)
}

// Debug logs message with Debug level
func Debug(msg string) {
	sugar.Debug(msg)
}

// Flush will flush zap buffers
func Flush() {
	sugar.Sync()
	logger.Sync()
}

func toZapLevel(str string) zapcore.Level {
	switch str {
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "debug":
		return zap.DebugLevel
	default:
		return zap.DebugLevel
	}
}
