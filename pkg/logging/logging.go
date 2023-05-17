package logging

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey int

const loggerKey contextKey = iota

var (
	defaultLogger     *zap.SugaredLogger
	defaultLoggerOnce sync.Once
)

var conf = &Config{
	Encoding:          "console",
	Level:             zapcore.InfoLevel,
	Development:       true,
	EncoderConfig:     NewEncoderConfig(),
	DisableStacktrace: true,
}

type Config struct {
	Encoding          string
	Level             zapcore.Level
	Development       bool
	EncoderConfig     zapcore.EncoderConfig
	DisableStacktrace bool
}

// SetConfig sets given logging configs for DefaultLogger's logger.
// Must set configs before calling DefaultLogger()
func SetConfig(c *Config) {
	conf = &Config{
		Encoding:          c.Encoding,
		Level:             c.Level,
		Development:       c.Development,
		EncoderConfig:     c.EncoderConfig,
		DisableStacktrace: c.DisableStacktrace,
	}
}

// NewLogger creates a new logger with the config.Context i.e config package should be initialized
func NewLogger() *zap.SugaredLogger {
	conf := zap.Config{
		Encoding:          conf.Encoding,
		EncoderConfig:     conf.EncoderConfig,
		Level:             zap.NewAtomicLevelAt(conf.Level),
		Development:       conf.Development,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: conf.DisableStacktrace,
	}
	logger, err := conf.Build()
	if err != nil {
		logger = zap.NewNop()
	}
	return logger.Sugar()
}

// DefaultLogger returns the default logger for the package.
func DefaultLogger() *zap.SugaredLogger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLogger()
	})
	return defaultLogger
}

// WithLogger creates a new context with the provided logger attached.
func WithLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext returns the logger stored in the context, otherwise a default logger is returned.
func FromContext(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return DefaultLogger()
	}
	if logger, ok := ctx.Value(loggerKey).(*zap.SugaredLogger); ok {
		return logger
	}
	return DefaultLogger()
}

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func GetNeme() string {
	return "Name~"
}
