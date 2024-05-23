package storage

import (
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

func PeriodicSend(c *websocket.Conn) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err := c.WriteMessage(websocket.TextMessage, []byte("Periodic message"))
		if err != nil {
			log.Println("Write periodic message error:", err)
			return
		}	
	}
}