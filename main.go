package main

import (
	"fmt"
	"log"
	"os"
	"visper/client"
	"visper/server"
)

const (
	HOST = "localhost"
	PORT = ":3333"
	TYPE = "tcp"
)

func main() {
	if len(os.Args) != 2 || (os.Args[1] != "server" && os.Args[1] != "client") {
		log.Fatal("Usage: visp [server|client]")
	}
	switch os.Args[1] {
	case "server":
		fmt.Println("Initializing server")
		server.Init()
	case "client":
		fmt.Println("Initializing client")
		client.Init()
	}

}
