package websockets

import (
	"torch/torch-sync/storage"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/contrib/websocket"
)

func WebsocketsMiddleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func SyncHandler(c *websocket.Conn) {
	userID := c.Query("userID")

	n := storage.NewNotifier()
	n.StartListening(c, "items_update__"+userID)
}
