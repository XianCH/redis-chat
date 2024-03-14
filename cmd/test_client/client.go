package main

// func main() {
//向服务端发送http连接（连接附带username属性），服务端将连接升级为websocket连接

//向连接发送message
// }

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	// 服务器地址
	serverURL := "ws://localhost:12346/ws"

	// 连接 WebSocket 服务器
	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()

	// 启动读取消息的协程
	go readMessages(conn)

	// 从终端读取消息并发送到服务器
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error sending message to server:", err)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("Error reading from stdin:", err)
	}
}

// 从 WebSocket 连接中读取并打印消息
func readMessages(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from server:", err)
			return
		}
		fmt.Println("Received message:", string(message))
	}
}
