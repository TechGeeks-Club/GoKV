package server

import (
	"fmt"
	"log"
	"net"

	"github.com/B-AJ-Amar/gokv/internal/store"
)

func RunServer() {
	memory := store.InMemoryStore{}
	port := 6379
	fmt.Println("Launching server...")
	fmt.Println("Listen on port")
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}

	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go HandleConnection(conn, &memory)
	}
}
