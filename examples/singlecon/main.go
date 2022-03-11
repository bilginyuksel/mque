package main

import (
	"log"
	"net"
)

func main() {
	con, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	info := make([]byte, 1024)
	intValue, err := con.Read(info)
	if err != nil {
		panic(err)
	}

	log.Println("int value:", intValue)
	log.Println("str value:", string(info))
}
