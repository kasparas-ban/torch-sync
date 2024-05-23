package websockets

import (
	"torch/torch-sync/storage"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/contrib/websocket"
)

const NOTIFY_CHANNEL_NAME = "items_update__4bax1usfu2uk"

func WebsocketsMiddleware(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}

	return fiber.ErrUpgradeRequired
}

func SyncHandler(c *websocket.Conn) {
	n := storage.NewNotifier()
	n.StartListening(c, NOTIFY_CHANNEL_NAME)
}
