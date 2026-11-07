package main

import (
	"log"
	"net"
	"sara/internal/handle"
	"sara/internal/storage"
)

func main() {
	store := storage.New()

	ln, err := net.Listen("tcp", ":4869")
	if err != nil {
		log.Fatal("Failed to open listener on :4869", err)
	}
	log.Println("Listener opened in :4869")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("Failed to accept connection: ", err)
			continue
		}
		log.Println("Accept new connection!")

		go handle.Handler(conn, store)
	}
}
