package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bilginyuksel/mque/internal/conn"
)

func main() {
	s := conn.Server{
		OnConnection: func(c conn.Conn) {
			log.Printf("accepted new connection with id: %v\n\n", c)
		},
		ConnectionWriteDeadline: 2 * time.Second,
		ConnectionReadDeadline:  2 * time.Second,
	}
	if err := s.Start(":8080"); err != nil {
		panic(err)
	}
	log.Println("started listening on port 8080")

	go s.Listen()
	defer s.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("shutting down")
}
