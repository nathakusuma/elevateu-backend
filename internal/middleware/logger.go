package middleware

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"

	"github.com/nathakusuma/elevateu-backend/pkg/log"
)

func LoggerConfig() fiber.Handler {
	config := fiberzerolog.Config{
		Logger:          log.NewLogger(),
		FieldsSnakeCase: true,
		Fields: []string{
			"referer",
			"ip",
			"host",
			"url",
			"ua",
			"latency",
			"status",
			"method",
		},
		Messages: []string{
			"[] Server error",
			"[] Client error",
			"[] Success",
		},
	}

	return fiberzerolog.New(config)
}
