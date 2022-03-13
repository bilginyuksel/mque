package broker

import (
	"log"

	"github.com/bilginyuksel/mque/internal/conn"
	"github.com/bilginyuksel/mque/internal/topic/v2"
)

type Writer struct {
	conn        *conn.Conn
	topicWriter topic.Writer
}

func NewWriter(c *conn.Conn, topicWriter topic.Writer) *Writer {
	return &Writer{
		conn:        c,
		topicWriter: topicWriter,
	}
}

func (w *Writer) WriteMessage() error {
	msg, err := w.conn.Read()
	if err != nil {
		return err
	}
	log.Println("msg:", string(msg))

	w.topicWriter.Write(msg)
	return nil
}
