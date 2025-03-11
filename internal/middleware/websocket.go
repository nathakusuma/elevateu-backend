package middleware

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func WebsocketUpgrade(ctx *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(ctx) {
		ctx.Locals("allowed", true)
		return ctx.Next()
	}
	return fiber.ErrUpgradeRequired
}
