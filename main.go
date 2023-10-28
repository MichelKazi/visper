package main

import (
	"fmt"
	"log"
	"os"
  "visper/config"
	"visper/client"
	"visper/server"
)

// TODO: read config from yaml

func main() {
	if len(os.Args) != 2 || (os.Args[1] != "server" && os.Args[1] != "client") {
		log.Fatal("Usage: visp [server|client]")
	}

  config, err := config.ReadConfig("visprc.yaml")
  if err != nil {
    log.Fatalf("Failed to read config: %v", err)
  }

  switch os.Args[1] {
  case "server":
    fmt.Println("Initializing server")
    server.Init()
  case "client":
    fmt.Println("Initializing client")
    host := config.Application.Server.Host
    port := config.Application.Server.Port
    protocol := config.Application.Server.Protocol
    serverAddr := fmt.Sprintf("%s:%d", host, port)

    cl := &client.Client{
      Protocol: protocol,
      ServerAddr:     serverAddr,
    }
    if err := cl.Init(); err != nil {
      panic(err)
    }
    cl.Start()

  }

}
