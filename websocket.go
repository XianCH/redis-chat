package redischat

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// websocket 升级
func RunSocket(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		log.Println("user name is empty")
		return
	}

	//升级服务
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	//create new client
	uuid := uuid.New()
	user := NewUser(uuid.String(), username)
	client := NewClient(conn, *user)
	MyServer.register <- client
	go client.Read()
	go client.Write()
	//todo：从redis中获取消息
	message, err := rdc.LRange(ctx, "group_chat_message", 0, -1).Result()
	if err != nil && err != redis.Nil {
		log.Printf("redis LRange error :%v", err)
		return
	}
	for _, msg := range message {
		data := []byte(msg)
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
