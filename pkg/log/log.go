package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/nathakusuma/elevateu-backend/domain/ctxkey"
	"github.com/nathakusuma/elevateu-backend/internal/infra/env"
)

var (
	logger zerolog.Logger
	once   sync.Once
)

func NewLogger() *zerolog.Logger {
	once.Do(func() {
		zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
			fn := runtime.FuncForPC(pc)
			if fn == nil {
				return fmt.Sprintf("%s:%d", filepath.Base(file), line)
			}
			return fmt.Sprintf("%s:%d %s", filepath.Base(file), line, filepath.Base(fn.Name()))
		}

		writers := []io.Writer{
			zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: time.RFC3339,
			},
		}

		// Add file output only if not in test environment
		if env.GetEnv().AppEnv != "test" {
			fileWriter := &lumberjack.Logger{
				Filename:   fmt.Sprintf("./storage/logs/app-%s.log", time.Now().Format("2006-01-02")),
				LocalTime:  true,
				Compress:   true,
				MaxSize:    100, // megabytes
				MaxAge:     7,   // days
				MaxBackups: 3,
			}
			writers = append(writers, fileWriter)
		}

		logger = zerolog.New(zerolog.MultiLevelWriter(writers...)).
			With().
			Timestamp().
			Logger()
	})

	return &logger
}

func logEvent(ctx context.Context, event *zerolog.Event, fields map[string]interface{}, msg string) {
	// Skip 2 frames: logEvent and the log level function (e.g., Info)
	event = event.Caller(2)

	// Extract user ID from context if it exists
	if ctx != nil {
		if userID := ctx.Value(ctxkey.UserID); userID != nil {
			if fields == nil {
				fields = make(map[string]interface{})
			}
			fields["requester.id"] = userID
		}
	}

	if len(fields) > 0 {
		entries := make([]interface{}, len(fields)*2)
		i := 0
		for k, v := range fields {
			entries[i] = k
			entries[i+1] = v
			i += 2
		}
		event.Fields(entries)
	}
	event.Msg(msg)
}

func Debug(ctx context.Context, fields map[string]interface{}, msg string) {
	if event := logger.Debug(); event.Enabled() {
		logEvent(ctx, event, fields, msg)
	}
}

func Info(ctx context.Context, fields map[string]interface{}, msg string) {
	if event := logger.Info(); event.Enabled() {
		logEvent(ctx, event, fields, msg)
	}
}

func Warn(ctx context.Context, fields map[string]interface{}, msg string) {
	if event := logger.Warn(); event.Enabled() {
		logEvent(ctx, event, fields, msg)
	}
}

func Error(ctx context.Context, fields map[string]interface{}, msg string) {
	if event := logger.Error(); event.Enabled() {
		logEvent(ctx, event, fields, msg)
	}
}

func ErrorWithTraceID(ctx context.Context, fields map[string]interface{}, msg string) uuid.UUID {
	traceID, err := uuid.NewRandom()
	if err != nil {
		Error(ctx, map[string]interface{}{
			"error": err.Error(),
		}, "[log.ErrorWithTraceID] failed to generate trace ID")
	}

	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["trace_id"] = traceID

	if event := logger.Error(); event.Enabled() {
		logEvent(ctx, event, fields, msg)
	}

	return traceID
}

func Fatal(ctx context.Context, fields map[string]interface{}, msg string) {
	logEvent(ctx, logger.Fatal(), fields, msg)
}

func Panic(ctx context.Context, fields map[string]interface{}, msg string) {
	logEvent(ctx, logger.Panic(), fields, msg)
}
