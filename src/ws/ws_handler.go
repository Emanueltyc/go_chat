package ws

import (
	"go_chat/src/types"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	const userContextKey types.ContextJWTClaimKey = "userID"

	userID, ok := r.Context().Value(userContextKey).(string)
	if !ok {
		return
	}

	client := &Client{
		UserID: userID,
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	hub.Register <- client

	go client.WritePump()
	client.ReadPump()

	defer func() {
		hub.Unregister <- client
		client.Conn.Close()
	}()
}
