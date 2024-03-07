package redischat

import (
	"log"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/x14n/redis-chat/pb"
)

type User struct {
	ID   string
	name string
}

type Client struct {
	conn *websocket.Conn
	user *User
	Send chan []byte
}

func NewUser(ID string, name string) *User {
	return &User{
		ID:   ID,
		name: name,
	}
}

func NewClient(conn *websocket.Conn, user User) *Client {
	return &Client{
		conn: conn,
		user: &user,
		Send: make(chan []byte),
	}
}

func (c *Client) Read() {
	defer func() {
		c.conn.Close()
		close(c.Send)
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("message read error:%v", err)
			MyServer.unregister <- c
			c.conn.Close()
			break
		}

		msg := &pb.Message{}
		err = proto.Unmarshal(message, msg)
		if err != nil {
			log.Printf("proto unmarshal error :%v", err)
		}
		MyServer.brocast <- message
	}
}

func (client *Client) Write() {
	defer client.conn.Close()
	for message := range client.Send {
		client.conn.WriteMessage(websocket.BinaryMessage, message)
	}
}
