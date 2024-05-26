package storage

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/lib/pq"
)

type notifier struct {
	userID string
	// DB connection
	listener *pq.Listener
	dbFailed chan error
	// Client connection
	clientFailed chan error
}

var MIN_RECONN = 10 * time.Second
var MAX_RECONN = 30 * time.Second

func NewNotifier(userID string) *notifier {
	l := pq.NewListener(DSN, MIN_RECONN, MAX_RECONN, handleListenerError)

	n := &notifier{
		userID:       userID,
		dbFailed:     make(chan error, 2),
		clientFailed: make(chan error, 2),
		listener:     l,
	}

	return n
}

func (n *notifier) StartListening(c *websocket.Conn, channelName string) {
	go n.handleClientConn(c)
	go n.handleDBConn(channelName)

	for {
		select {
		case e := <-n.listener.Notify:
			err := c.WriteMessage(1, []byte(e.Extra))
			if err != nil {
				n.clientFailed <- err
			}
		case err := <-n.dbFailed:
			n.listener.Close()
			fmt.Printf("DB connection error: %v \n", err)
			return
		case err := <-n.clientFailed:
			n.listener.Close()
			fmt.Printf("Client connection error: %v \n", err)
			return
		}
	}
}

func (n *notifier) handleDBConn(channelName string) {
	err := n.listener.Listen(channelName)
	if err != nil {
		n.dbFailed <- err
	}
}

func (n *notifier) handleClientConn(c *websocket.Conn) {
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			n.clientFailed <- err
			break
		}

		err = ProcessCmd(msg, n.userID)
		if err != nil {
			log.Printf("Command failed: %v\n", err)
		}
	}
}
