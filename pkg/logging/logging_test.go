package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	assert.NotNil(t, NewLogger())
}

func TestDefaultLogger(t *testing.T) {
	l1 := DefaultLogger()
	assert.NotNil(t, l1)

	l2 := DefaultLogger()
	assert.NotNil(t, l2)

	assert.Equal(t, l1, l2)
}

func TestFromContext(t *testing.T) {
	cases := []struct {
		name string
		ctx  context.Context
	}{
		{name: "DefaultContext", ctx: context.Background()},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.ctx

			l1 := FromContext(ctx)
			assert.NotNil(t, l1)

			ctx = WithLogger(ctx, l1)
			l2 := FromContext(ctx)

			// then
			assert.Equal(t, l1, l2)
		})
	}
}

func TestFromContextWithNil(t *testing.T) {
	l1 := DefaultLogger()

	l2 := FromContext(nil)

	assert.Equal(t, l1, l2)
}

func TestLoggingFormat(t *testing.T) {
	SetConfig(&Config{Encoding: "json", Level: -1})

	output := captureLoggingOutput(func() {
		DefaultLogger().Errorw("my log message", "x-request-id", "request1")
	})
	t.Log(output)

	var res map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(output), &res))

	assert.Equal(t, "ERROR", res["L"].(string))
	assert.Equal(t, "my log message", res["M"].(string))
	assert.Contains(t, res["S"].(string), "pkg/logging.TestLoggingFormat")
	parsedTime, err := time.Parse("2006-01-02T15:04:05.000Z0700", res["T"].(string))
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Now(), parsedTime, time.Hour)
	assert.Contains(t, res["C"].(string), "logging/logging_test.go")
	assert.Equal(t, "request1", res["x-request-id"].(string))
}

func captureLoggingOutput(doFunc func()) string {
	DefaultLogger() // call to initialize default logger if not exist
	oldLogger := defaultLogger
	sink := &MemorySink{new(bytes.Buffer)}
	zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
		return sink, nil
	})

	cfg := zap.Config{
		Encoding:         conf.Encoding,
		EncoderConfig:    NewEncoderConfig(),
		Level:            zap.NewAtomicLevelAt(conf.Level),
		Development:      conf.Development,
		OutputPaths:      []string{"memory://"},
		ErrorOutputPaths: []string{"memory://"},
	}

	newLogger, _ := cfg.Build()
	defaultLogger = newLogger.Sugar()

	doFunc()

	defaultLogger = oldLogger
	return sink.String()
}

// MemorySink implements zap.Sink to collect all messages in buffer.
type MemorySink struct {
	*bytes.Buffer
}

func (s *MemorySink) Close() error { return nil }
func (s *MemorySink) Sync() error  { return nil }
