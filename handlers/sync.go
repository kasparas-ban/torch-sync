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

	wsId, ok := c.Locals("ws_id").(string)
	if !ok || wsId == "" {
		c.Close()
		return
	}

	n := storage.NewNotifier(userID, wsId)
	n.StartListening(c, "items_update__"+userID)
}
