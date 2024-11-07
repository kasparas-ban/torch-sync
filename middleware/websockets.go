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
	err := VerifyUpdateToken(c, token)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	// Add websocket ID to the context
	wsId := strings.Trim(data[1], " ")
	c.Locals("ws_id", wsId)

	if websocket.IsWebSocketUpgrade(c) {
		// Handle the WebSocket handshake to set the subprotocol
		secWebSocketProtocol := c.Get("Sec-WebSocket-Protocol")
		selectedProtocol := ""
		if secWebSocketProtocol != "" {
			protocols := strings.Split(secWebSocketProtocol, ",")
			selectedProtocol = protocols[1]
		}

		c.Context().Response.Header.Set("Sec-WebSocket-Protocol", selectedProtocol)

		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}
