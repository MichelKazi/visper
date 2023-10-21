package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	HOST = "127.0.0.1"
	PORT = ":8080"
	TYPE = "tcp"
)

func Init() {
	conn, err := net.Dial(TYPE, HOST+PORT)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter your name: ")
	scanner.Scan()
	name := scanner.Text()
	conn.Write([]byte(name + " has joined the chat!\n"))

	defer conn.Close()
	go readMessages(conn)

	for scanner.Scan() {
    h, m, s := time.Now().Clock()
    timestamp := fmt.Sprintf("%d:%d:%d", h, m, s)
    msg := fmt.Sprintf("[%s] %s: %s", timestamp, name, scanner.Text() + "\n")
		conn.Write([]byte(msg))
	}
}

func readMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Lost connection to server.")
		}
		fmt.Println("Received:", msg)
	}
}
