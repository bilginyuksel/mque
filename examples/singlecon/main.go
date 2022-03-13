package main

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
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
		Type:         Subscriber,
	}
	confBytes, _ := json.Marshal(conf)
	if _, err := con.Write(confBytes); err != nil {
		panic(err)
	}

	go func() {
		chunk := make([]byte, 2048)
		for {
			msgByteLength, err := con.Read(chunk)
			if err != nil {
				log.Printf("msg read failed, err: %v\n", err)
				break
			}

			msg := make([]byte, msgByteLength)
			copy(msg, chunk)

			log.Printf("length: %d, msg: %s\n", msgByteLength, string(msg))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("shutting down")

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
