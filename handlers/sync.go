package handlers

import (
	"torch/torch-sync/middleware"
	"torch/torch-sync/storage"

	"github.com/gofiber/contrib/websocket"
)

func SyncHandler(c *websocket.Conn) {
	userID, ok := c.Locals(middleware.UserID_metadata).(string)
	if !ok || userID == "" {
		c.Close()
		return
	}

	n := storage.NewNotifier()
	n.StartListening(c, "items_update__"+userID)
}
