package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Cors() fiber.Handler {
	config := cors.Config{
		AllowMethods:  "GET,POST,PUT,DELETE,PATCH,OPTIONS,HEAD",
		AllowHeaders:  "Content-Type,Authorization,Accept,Origin,X-Requested-With,X-XSRF-Token,X-Cursor,Token-Type",
		ExposeHeaders: "Content-Length",
	}

	return cors.New(config)
}
