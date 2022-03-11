package main

import (
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	log.Println("started listening on port 8080")

	conn, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	log.Println("accepted connection")

	intValue, err := conn.Write([]byte("hello"))
	if err != nil {
		panic(err)
	}

	log.Println("int value:", intValue)
}
