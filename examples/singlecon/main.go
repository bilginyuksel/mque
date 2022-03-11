package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/google/uuid"
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
		Topic:        fmt.Sprintf("test-topic-%s", uuid.NewString()),
		ReadTimeout:  1000,
		WriteTimeout: 2000,
		MaxByteSize:  33223232,
		MinByteSize:  10,
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

			log.Println("received message from server:", string(msg))
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

type Config struct {
	Topic        string `json:"topic"`
	ReadTimeout  int64  `json:"read_timeout"`
	WriteTimeout int64  `json:"write_timeout"`
	MaxByteSize  int64  `json:"max_byte_size"`
	MinByteSize  int64  `json:"min_byte_size"`

	// Encryption
}
