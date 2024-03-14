package main

import redischat "github.com/x14n/redis-chat"

func main() {
	server := redischat.NewServer()
	server.Score()
	server.Start()
}
