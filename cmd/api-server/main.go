package main

import (
	"log"
	"net"

	"github.com/bmolinari/go-http-server/internal/routes"
)

func main() {
	server := routes.NewRouter()

	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Println("Server is running on port 2000...")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go server.HandleConnection(conn)
	}
}
