package server

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"slices"
	"strings"
)

type client struct {
	username string
	conn     net.Conn
}

var (
	pool []client
	motd string
)

func Run(address string) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("failed to listen to connection:", err)
		return
	}
	defer ln.Close()
	log.Println("listening on", address)
	done := make(chan bool)
	go accept(ln, done)
	go write(done)
	<-done
	log.Println("shutting down")
}

func accept(ln net.Listener, done chan bool) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("failed to accept connection:", err)
			return
		}
		go read(conn, done)
	}
}

func read(conn net.Conn, done chan bool) {
	username := username(conn)
	client := client{username: username, conn: conn}
	pool = append(pool, client)
	conn.Write([]byte(motd))
	fmt.Print("\r")
	log.Println("new connection:", conn.RemoteAddr().String(), ":", username)
	fmt.Print("> ")
	scanner := bufio.NewScanner(conn)
	for {
		select {
		case <-done:
			goto done
		default:
			if scanner.Scan() {
				msg := scanner.Text()
				words := strings.Split(msg, " ")
				switch words[0] {
				case "exit":
					goto done
				case "help":
					conn.Write([]byte(
						"\rexit - close the app\n\r" +
							"help - display this message\n\r" +
							"list - display all connected clients\n\r" +
							"msg <username> <message> - send private message to the user\n",
					))
				case "msg":
					msg = fmt.Sprintf("\033[2K%s [private]: %s\n", client.username, strings.Join(words[2:], " "))
					message(words[1], msg)
				case "list":
					fmt.Print("\r")
					for _, client := range pool {
						conn.Write([]byte("\r" + client.username + "\n"))
					}
				default:
					msg = fmt.Sprintf("\033[2K%s: %s\n", client.username, msg)
					broadcast(conn, msg)
				}
			} else {
				goto done
			}
		}
	}
done:
	pool = remove(pool, client)
	conn.Close()
	fmt.Print("\r")
	log.Println("disconnected:", conn.RemoteAddr().String())
	fmt.Print("> ")
}

func broadcast(cur net.Conn, msg string) {
	for _, client := range pool {
		if client.conn == cur {
			continue
		}
		client.conn.Write([]byte(msg))
	}
	fmt.Print(msg)
	fmt.Print("> ")
}

func remove[Type comparable](s []Type, v Type) []Type {
	for i, e := range s {
		if v == e {
			return slices.Delete(s, i, i+1)
		}
	}
	return s
}

func write(done chan bool) {
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		words := strings.Split(msg, " ")
		switch words[0] {
		case "exit":
			goto done
		case "help":
			fmt.Print(
				"exit - close the app\n" +
					"help - display this message\n" +
					"kick <username> - kick user from chat\n" +
					"list - display all connected clients\n" +
					"motd <message> - set message of the day\n" +
					"msg <username> <message> - send private message to the user\n",
			)
		case "kick":
			kick(words[1])
		case "list":
			fmt.Print("\r")
			for _, client := range pool {
				fmt.Println(client.username)
			}
		case "motd":
			motd = fmt.Sprintf("\033[2Kmotd: %s\n",
				strings.Join(words[1:], " "))
		case "msg":
			user := words[1]
			msg := fmt.Sprintf("\033[2Kadmin [private]: %s\n",
				strings.Join(words[2:], " "))
			message(user, msg)
		default:
			msg := fmt.Sprintf("\033[2Kadmin: %s\n", msg)
			for _, client := range pool {
				client.conn.Write([]byte(msg))
			}
		}
		fmt.Print("> ")
	}
done:
	done <- true
}

func kick(username string) {
	for _, client := range pool {
		if client.username == username {
			client.conn.Close()
			pool = remove(pool, client)
			return
		}
	}
}

func message(user, message string) {
	if user == "admin" {
		fmt.Print(message)
		fmt.Print("> ")
		return
	}
	for _, client := range pool {
		if client.username == user {
			client.conn.Write([]byte(message))
			return
		}
	}
}

func username(conn net.Conn) string {
	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Println("failed to read from connection:", err)
		return "user"
	}
	name = strings.TrimSpace(name)
	bytes := make([]byte, 2)
	rand.Read(bytes)
	hexChars := hex.EncodeToString(bytes)
	name += "#" + hexChars
	return name
}
