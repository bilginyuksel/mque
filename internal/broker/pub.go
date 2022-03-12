package broker

import (
	"log"

	"github.com/bilginyuksel/mque/internal/conn"
	"github.com/bilginyuksel/mque/internal/topic"
)

type ReaderWriter struct {
	conn        *conn.Conn
	offset      int
	topicReader topic.Reader
	topicWriter topic.Writer
}

func NewReaderWriter(conn *conn.Conn, topic topic.Topic) *ReaderWriter {
	return &ReaderWriter{
		conn:        conn,
		topicReader: topic,
		topicWriter: topic,
	}
}

func (rw *ReaderWriter) ReadMessage() error {
	// Check if offset is bigger than or equal to queue length
	// if it is then we need to wait for a new element to be inserted to queue
	if rw.offset >= rw.topicReader.Size() {
		<-rw.topicReader.Listen()
	}

	msg := rw.topicReader.Get(rw.offset)
	if err := rw.conn.Write(msg); err != nil {
		return err
	}

	rw.offset++
	return nil
}

func (rw *ReaderWriter) WriteMessage() error {
	msg, err := rw.conn.Read()
	if err != nil {
		return err
	}
	log.Println("received message:", msg)

	return rw.topicWriter.Push(msg)
}
