package broker

import (
	"log"

	"github.com/bilginyuksel/mque/internal/conn"
	"github.com/bilginyuksel/mque/internal/topic/v2"
)

type Broker struct {
	topics map[string]topic.Topic
}

func New() *Broker {
	return &Broker{
		topics: make(map[string]topic.Topic),
	}
}

func (b *Broker) createIfNotExists(topicName string) topic.Topic {
	if t, ok := b.topics[topicName]; ok {
		return t
	}

	t, err := topic.New(topicName)
	if err != nil {
		log.Println("error while creating topic:", err)
		return nil
	}

	b.topics[topicName] = t
	return t
}

func (b *Broker) Subscribe(topicName string, connection *conn.Conn) {
	t := b.createIfNotExists(topicName)
	reader := NewReader(connection, t.CreateReader())

	go func() {
		for {
			if err := reader.ReadMessage(); err != nil {
				log.Println("error while reading message:", err)
				break
			}
		}
	}()
}

func (b *Broker) Publish(topicName string, connection *conn.Conn) {
	t := b.createIfNotExists(topicName)
	writer := NewWriter(connection, t.CreateWriter())

	go func() {
		for {
			if err := writer.WriteMessage(); err != nil {
				log.Println("error while writing message:", err)
				break
			}
		}
	}()
}
