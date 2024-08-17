package handlers

import (
	"log/slog"
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

	wsID, ok := c.Locals("ws_id").(string)
	if !ok || wsID == "" {
		c.Close()
		return
	}

	n := storage.NewNotifier(userID, wsID)
	slog.Info("Listening on DB channel", "wsID", wsID, "userID", userID)
	n.StartListening(c, "db_update__"+userID)
}
