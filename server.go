package redischat

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	logger "github.com/x14n/redis-chat/log"
	"github.com/x14n/redis-chat/pb"
)

var xlog = logger.GetLogger

type server struct {
	Clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	brocast    chan []byte
	mu         sync.Mutex
}

func NewServer() *server {
	return &server{
		Clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		brocast:    make(chan []byte),
	}
}

var MyServer = NewServer()

func (s *server) Close() {
	for key := range s.Clients {
		s.Clients[key].conn.Close()
		delete(s.Clients, key)
	}
	close(s.register)
	close(s.unregister)
}

// 开启注册，取消注册，广播通道
func (s *server) Start() {
	for {
		select {
		case conn := <-s.register:
			s.Clients[conn.user.name] = conn
			log.Printf("one client register:%s", conn.user.name)
			message := &pb.Message{
				From:    "system",
				Content: "welecom",
			}
			proto, err := proto.Marshal(message)
			if err != nil {
				log.Printf("proto marshal error:%v", err)
			}
			conn.Send <- proto

		case conn := <-s.unregister:
			log.Printf("one client close the connection:%v", conn.user.name)
			if _, ok := s.Clients[conn.user.name]; ok {
				delete(s.Clients, conn.user.name)
				close(conn.Send)
			}
		case message := <-s.brocast:
			log.Printf("brocast to every client")
			for _, client := range s.Clients {
				client.Send <- message
			}
		}
	}
}

func (s *server) Score() {
	router := http.NewServeMux()
	router.HandleFunc("/ws", RunSocket)
	s1 := &http.Server{
		Addr:           ":12346",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("server start at 12346")
	err := s1.ListenAndServe()
	if err != nil {
		log.Printf("server start error :%v\n", err)
		return
	}
}
