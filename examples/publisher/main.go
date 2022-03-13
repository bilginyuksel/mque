package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
)

func main() {
	con, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	info := make([]byte, 1024)
	if _, err := con.Read(info); err != nil {
		panic(err)
	}

	log.Println("received value from server:", string(info))

	conf := Config{
		Topic:        "test",
		ReadTimeout:  1000,
		WriteTimeout: 2000,
		MaxByteSize:  33223232,
		MinByteSize:  1,
		Type:         Publisher,
	}
	confBytes, _ := json.Marshal(conf)
	if _, err := con.Write(confBytes); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		var msgContent string
		if scanner.Scan() {
			msgContent = scanner.Text()
		}

		if msgContent == "q" {
			break
		}

		log.Println("sending message:", msgContent)

		if _, err := con.Write([]byte(msgContent)); err != nil {
			log.Println("write failed, err:", err)
		}
	}

	if err := con.Close(); err != nil {
		panic(err)
	}
}

const (
	Publisher = iota + 1
	Subscriber
)

type Config struct {
	Topic        string `json:"topic"`
	ReadTimeout  int64  `json:"read_timeout"`
	WriteTimeout int64  `json:"write_timeout"`
	MaxByteSize  int64  `json:"max_byte_size"`
	MinByteSize  int64  `json:"min_byte_size"`
	Type         uint8  `json:"type"`

	// Encryption
}
