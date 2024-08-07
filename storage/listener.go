package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/lib/pq"
	"github.com/tidwall/gjson"
)

type notifier struct {
	userID string
	wsId   string
	// DB connection
	listener *pq.Listener
	dbFailed chan error
	// Client connection
	clientFailed chan error
}

var MIN_RECONN = 10 * time.Second
var MAX_RECONN = 30 * time.Second

func NewNotifier(userID string, wsId string) *notifier {
	l := pq.NewListener(DSN, MIN_RECONN, MAX_RECONN, handleListenerError)

	n := &notifier{
		userID:       userID,
		wsId:         wsId,
		dbFailed:     make(chan error, 2),
		clientFailed: make(chan error, 2),
		listener:     l,
	}

	return n
}

func (n *notifier) StartListening(c *websocket.Conn, channelName string) {
	go safeGo(func() { n.handleClientConn(c) })
	n.handleDBConn(channelName)

	for {
		select {
		case e := <-n.listener.Notify:
			wsId, op, err := readMsg(e.Extra)
			if err != nil || (wsId == n.wsId && op != "UPDATE") {
				continue
			}

			err = c.WriteMessage(1, []byte(e.Extra))
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

		err = ProcessCmd(msg, n.userID, n.wsId)
		if err != nil {
			log.Printf("Command failed: %v\n", err)
		}
	}
}

func safeGo(fn func()) {
	go func() {
		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Recovered from panic:", r)
					}
				}()
				fn()
			}()
		}
	}()
}

func AddWebsocketIdJSON(msg string, wsId string) (string, error) {
	var data map[string]json.RawMessage
	if err := json.Unmarshal([]byte(msg), &data); err != nil {
		return "", err
	}

	newKey := "wsId"
	data[newKey] = json.RawMessage(`"` + wsId + `"`)

	updatedJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(updatedJSON), nil
}

func readMsg(msg string) (string, string, error) {
	wsIdRes := gjson.GetBytes([]byte(msg), "ws_id")
	opRes := gjson.GetBytes([]byte(msg), "op")

	wsId := wsIdRes.String()
	op := opRes.String()

	if wsId == "" || op == "" {
		return "", "", errors.New("failed to read notification message")
	}

	return wsId, op, nil
}
