package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/contrib/websocket"
)

func WebsocketsMiddleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}
