package storage

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/lib/pq"
	"github.com/tidwall/gjson"
	"golang.org/x/exp/slog"
)

type notifier struct {
	userID string
	wsID   string
	// DB connection
	listener *pq.Listener
	dbFailed chan error
	// Client connection
	clientFailed chan error
}

var MIN_RECONN = 10 * time.Second
var MAX_RECONN = 30 * time.Second

func NewNotifier(userID string, wsID string) *notifier {
	l := pq.NewListener(DSN, MIN_RECONN, MAX_RECONN, handleListenerError)

	n := &notifier{
		userID:       userID,
		wsID:         wsID,
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
			slog.Info("Msg from DB", "msg", e.Extra)
			wsID, op, err := readMsg(e.Extra)
			if err != nil || (wsID == n.wsID && op != "UPDATE") {
				continue
			}

			err = c.WriteMessage(1, []byte(e.Extra))
			if err != nil {
				n.clientFailed <- err
			}
		case err := <-n.dbFailed:
			n.listener.Close()
			slog.Error("DB connection error", "error", err)
			return
		case err := <-n.clientFailed:
			n.listener.Close()
			slog.Error("Client connection error", "error", err)
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

		err = ProcessCmd(msg, n.userID, n.wsID)
		if err != nil {
			slog.Error("Command failed", "err", err)
		}
	}
}

func safeGo(fn func()) {
	go func() {
		for {
			func() {
				defer func() {
					if r := recover(); r != nil {
						slog.Error("Recovered from panic:", "data", r)
					}
				}()
				fn()
			}()
		}
	}()
}

func AddWebsocketIdJSON(msg string, wsID string) (string, error) {
	var data map[string]json.RawMessage
	if err := json.Unmarshal([]byte(msg), &data); err != nil {
		return "", err
	}

	newKey := "wsID"
	data[newKey] = json.RawMessage(`"` + wsID + `"`)

	updatedJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(updatedJSON), nil
}

func readMsg(msg string) (string, string, error) {
	wsIDRes := gjson.GetBytes([]byte(msg), "ws_id")
	opRes := gjson.GetBytes([]byte(msg), "op")

	wsID := wsIDRes.String()
	op := opRes.String()

	if wsID == "" || op == "" {
		return "", "", errors.New("failed to read notification message")
	}

	return wsID, op, nil
}
