package redischat

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunSocket(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("user")
	if user == "" {
		return
	}
	ws, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	id := string(uuid.New())
	user := &User{
        ID:id,
		name: user,
	}
	client := &Client{
		Send: make(chan []byte),
	}

	MyServer.Register <- client
	go client.Read()
	go client.write()
}
