package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/nathakusuma/elevateu-backend/internal/infra/env"
)

var (
	logger zerolog.Logger
	once   sync.Once
)

func NewLogger() *zerolog.Logger {
	once.Do(func() {
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

func logEvent(event *zerolog.Event, fields map[string]interface{}, msg string) {
	if len(fields) > 0 {
		entries := make([]interface{}, 0, len(fields)*2)
		for k, v := range fields {
			entries = append(entries, k, v)
		}
		event.Fields(entries)
	}
	event.Msg(msg)
}

func Debug(fields map[string]interface{}, msg string) {
	if event := logger.Debug(); event.Enabled() {
		logEvent(event, fields, msg)
	}
}

func Info(fields map[string]interface{}, msg string) {
	if event := logger.Info(); event.Enabled() {
		logEvent(event, fields, msg)
	}
}

func Warn(fields map[string]interface{}, msg string) {
	if event := logger.Warn(); event.Enabled() {
		logEvent(event, fields, msg)
	}
}

func Error(fields map[string]interface{}, msg string) {
	if event := logger.Error(); event.Enabled() {
		logEvent(event, fields, msg)
	}
}

func ErrorWithTraceID(fields map[string]interface{}, msg string) uuid.UUID {
	traceID, err := uuid.NewRandom()
	if err != nil {
		Error(map[string]interface{}{
			"error": err.Error(),
		}, "[log.ErrorWithTraceID] failed to generate trace ID")
	}

	if event := logger.Error(); event.Enabled() {
		fields["trace_id"] = traceID
		logEvent(event, fields, msg)
	}

	return traceID
}

func Fatal(fields map[string]interface{}, msg string) {
	logEvent(logger.Fatal(), fields, msg)
}

func Panic(fields map[string]interface{}, msg string) {
	logEvent(logger.Panic(), fields, msg)
}
