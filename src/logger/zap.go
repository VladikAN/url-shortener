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

// Fatalf logs formatted message with Fatal level
func Fatalf(msg string, arg ...interface{}) {
	sugar.Fatalf(msg, arg)
}

// Info logs message with Info level
func Info(msg string) {
	sugar.Info(msg)
}

// Infof logs formatted message with Info level
func Infof(msg string, arg ...interface{}) {
	sugar.Infof(msg, arg)
}

// Warn logs message with Warn level
func Warn(msg string) {
	sugar.Warn(msg)
}

// Warnf logs formatted message with Warn level
func Warnf(msg string, arg ...interface{}) {
	sugar.Warnf(msg, arg)
}

// Debug logs message with Debug level
func Debug(msg string) {
	sugar.Debug(msg)
}

// Debugf logs formatted message with Debug level
func Debugf(msg string, arg ...interface{}) {
	sugar.Debugf(msg, arg)
}

// Debugw logs enriched message with Debug level
func Debugw(msg string, fields ...interface{}) {
	sugar.Debugw(msg, fields...)
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
