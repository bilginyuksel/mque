package broker

import (
	"github.com/bilginyuksel/mque/internal/conn"
	"github.com/bilginyuksel/mque/internal/topic/v2"
)

// TODO: Create consumer groups
// TODO: Start from the end of the topic
type Reader struct {
	conn              *conn.Conn
	offset            int
	writtenBufferSize int64
	topicReader       topic.Reader
}

func NewReader(c *conn.Conn, topicReader topic.Reader) *Reader {
	return &Reader{
		conn:        c,
		topicReader: topicReader,
	}
}

func (r *Reader) ReadMessage() error {
	msg, err := r.topicReader.Read(r.writtenBufferSize, r.offset)
	if err != nil {
		return err
	}

	if err := r.conn.Write(msg); err != nil {
		return err
	}

	r.offset++
	r.writtenBufferSize += int64(len(msg)) + 11
	return nil
}
