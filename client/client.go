package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func Run(address, username string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Println("failed to connect to the server:", err)
		return
	}
	defer conn.Close()
	log.Println("connected to", address)
	conn.Write([]byte(username + "\n"))
	done := make(chan bool)
	go read(conn, done)
	go write(conn, done)
	<-done
	log.Println("disconnected")
}

func read(conn net.Conn, done chan bool) {
	scanner := bufio.NewScanner(conn)
	for {
		select {
		case <-done:
			return
		default:
			if scanner.Scan() {
				msg := scanner.Text()
				fmt.Println(msg)
				fmt.Print("> ")
			} else {
				done <- true
				return
			}
		}
	}
}

func write(conn net.Conn, done chan bool) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		input := scanner.Text() + "\n"
		_, err := conn.Write([]byte(input))
		if err != nil {
			log.Println("failed to write to connection:", err)
			return
		}
		fmt.Print("> ")
	}
	done <- true
}
