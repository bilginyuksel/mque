package broker

import (
	"log"

	"github.com/bilginyuksel/mque/internal/conn"
	"github.com/bilginyuksel/mque/internal/topic"
)

type Broker struct {
	topics map[string]topic.Topic
}

func New() *Broker {
	return &Broker{
		topics: make(map[string]topic.Topic),
	}
}

func (b *Broker) getOrCreate(topicName string) topic.Topic {
	if t, ok := b.topics[topicName]; ok {
		return t
	}

	t := topic.New(topic.NoAck)
	b.topics[topicName] = t
	return t
}

func (b *Broker) Subscribe(topicName string, connection *conn.Conn) {
	t := b.getOrCreate(topicName)
	readerWriter := NewReaderWriter(connection, t)

	go func() {
		for {
			if err := readerWriter.ReadMessage(); err != nil {
				log.Println("error while reading message:", err)
				break
			}
		}
	}()
}

func (b *Broker) Publish(topicName string, connection *conn.Conn) {
	t := b.getOrCreate(topicName)
	readerWriter := NewReaderWriter(connection, t)

	go func() {
		for {
			if err := readerWriter.WriteMessage(); err != nil {
				log.Println("error while writing message:", err)
				break
			}
		}
	}()
}
