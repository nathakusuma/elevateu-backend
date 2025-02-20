package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
)

func Compress() fiber.Handler {
	config := compress.Config{
		Next:  nil,
		Level: compress.LevelDefault,
	}

	return compress.New(config)
}
