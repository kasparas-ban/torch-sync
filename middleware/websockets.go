package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/contrib/websocket"
)

func WebsocketsMiddleware(c *fiber.Ctx) error {
	allHeaders := c.GetReqHeaders()
	protocols := allHeaders["Sec-Websocket-Protocol"]
	if len(protocols) != 1 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	token := protocols[0]
	err := VerifyToken(c, token)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}
