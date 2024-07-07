package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/contrib/websocket"
)

func WebsocketsMiddleware(c *fiber.Ctx) error {
	// Read protocol header
	allHeaders := c.GetReqHeaders()
	protocols := allHeaders["Sec-Websocket-Protocol"]
	data := strings.Split(protocols[0], ",")
	if len(data) != 2 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Verify auth token
	token := data[0]
	err := VerifyToken(c, token)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Add websocket ID to the context
	wsId := data[1]
	c.Locals("ws_id", wsId)

	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}
