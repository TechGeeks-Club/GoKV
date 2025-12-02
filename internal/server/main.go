package server

import (
	"fmt"
	"net"
)

func main() {
	port := 8080
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
		go HandleConnection(conn)
	}
}