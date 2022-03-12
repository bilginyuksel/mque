package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bilginyuksel/mque/internal/broker"
	"github.com/bilginyuksel/mque/internal/conn"
)

func main() {
	b := broker.New()
	s := conn.Server{
		OnConnection: func(c *conn.Conn) {
			// log.Printf("accepted new connection with id: %v\n\n", c)
			switch c.Conf.Type {
			case conn.Publisher:
				b.Publish(c.Conf.Topic, c)
				log.Println("created a publisher")
			case conn.Subscriber:
				b.Subscribe(c.Conf.Topic, c)
				log.Println("created a subscriber")
			}
		},
		Connections:             make(map[string]*conn.Conn),
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
