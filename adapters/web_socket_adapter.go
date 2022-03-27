package adapters

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	webSocketConnectionOnce     sync.Once
	webSocketConnectionInstance *webSocketConnection
)

func WebSocketAdapter() *webSocketConnection {
	webSocketConnectionOnce.Do(func() {
		webSocketConnectionInstance = &webSocketConnection{
			connections: map[string]*websocket.Conn{},
		}
	})

	return webSocketConnectionInstance
}

type webSocketConnection struct {
	connections map[string]*websocket.Conn
}

func (w *webSocketConnection) AddConnection(webSocketId string, conn *websocket.Conn) {
	if w.connections == nil {
		w.connections = map[string]*websocket.Conn{}
	}

	w.connections[webSocketId] = conn

	go func() {
		defer func() {
			if err := conn.Close(); err != nil {
				log.Error("web socket close", err)
			}
			delete(w.connections, webSocketId)
		}()
		for {
			if err := conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
				log.Error("web socket write", err)
				break
			}

			msgType, bytes, err := conn.ReadMessage()
			if err != nil {
				log.Error("web socket closed.", err)
				return
			}

			// We don't recognize any message that is not "pong".
			if msg := string(bytes[:]); msgType != websocket.TextMessage && msg != "pong" {
				log.Error("web socket Unrecognized message received.", err)
				continue
			} else {
				log.Info("Received: pong.")
			}

			time.Sleep(5 * time.Second)
		}
	}()
}

func (w webSocketConnection) SendMessage(webSocketId string, msg interface{}) error {
	if w.connections[webSocketId] == nil {
		return errors.New("invalid socket")
	}

	if err := w.connections[webSocketId].WriteJSON(msg); err != nil {
		return errors.Wrap(err, "websocket send message error")
	}

	return nil
}

func (w webSocketConnection) BroadcastMessage(msg interface{}) error {
	for _, conn := range w.connections {
		if err := conn.WriteJSON(msg); err != nil {
			return errors.Wrap(err, "websocket BroadcastMessage error")
		}
	}

	return nil
}
