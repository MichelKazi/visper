package server

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

const (
	PORT      = ":8080"
	CONN_TYPE = "tcp"
)

var clients = make(map[*net.Conn]bool)
var mtx sync.Mutex

func Init() {
	listener, err := net.Listen(CONN_TYPE, PORT)
	if err != nil {
		panic(err)
	}

	defer listener.Close()
  fmt.Printf("%s connection listening on port %s", CONN_TYPE, PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		// Locking client map to prevent other goroutines from writing at the same time
		mtx.Lock()
		clients[&conn] = true
		mtx.Unlock()

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			mtx.Lock()
			delete(clients, &conn)
			mtx.Unlock()
			break
		}
    broadcast(msg, conn)
	}
}

func broadcast(msg string, origin net.Conn) {
  mtx.Lock()
  defer mtx.Unlock()

  for client := range clients {
    if *client != origin {
      (*client).Write([]byte(msg))
    }
  }
}
