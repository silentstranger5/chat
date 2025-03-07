package main

import (
	"chat/client"
	"chat/server"
	"flag"
	"log"
)

func main() {
	mode := flag.String("mode", "client", "usage mode: client, server")
	address := flag.String("address", "localhost:8080", "connection address")
	user := flag.String("user", "user", "your username")
	flag.Parse()
	if *mode == "client" {
		client.Run(*address, *user)
	} else if *mode == "server" {
		server.Run(*address)
	} else {
		log.Println("invalid usage mode:", *mode)
	}
}
