package broker

import (
	"github.com/bilginyuksel/mque/internal/conn"
	"github.com/bilginyuksel/mque/internal/topic/v2"
)

type Reader struct {
	conn        *conn.Conn
	offset      int
	topicReader topic.Reader
}

func NewReader(c *conn.Conn, topicReader topic.Reader) *Reader {
	return &Reader{
		conn:        c,
		topicReader: topicReader,
	}
}

func (r *Reader) ReadMessage() error {
	msg := r.topicReader.Read(r.offset)
	if err := r.conn.Write(msg); err != nil {
		return err
	}

	r.offset++
	return nil
}
