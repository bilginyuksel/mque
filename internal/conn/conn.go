package conn

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
)

var (
	ErrByteLengthTooSmall = errors.New("byte length too small")
)

const (
	Publisher = iota + 1
	Subscriber
)

type Reader interface {
	Read() ([]byte, error)
}

type Writer interface {
	Write([]byte) error
}

type (
	Config struct {
		Topic        string `json:"topic"`
		ReadTimeout  int64  `json:"read_timeout"`
		WriteTimeout int64  `json:"write_timeout"`
		MaxByteSize  int64  `json:"max_byte_size"`
		MinByteSize  int64  `json:"min_byte_size"`
		Acks         uint8  `json:"acks"`
		Type         uint8  `json:"type"`

		// Encryption
	}

	Conn struct {
		conn net.Conn

		ID   string
		At   time.Time
		Conf Config

		// chunk use for communication
		chunk []byte
	}
)

func New(conn net.Conn, conf Config) (*Conn, error) {
	connection := &Conn{
		conn:  conn,
		ID:    uuid.NewString(),
		At:    time.Now(),
		Conf:  conf,
		chunk: make([]byte, conf.MaxByteSize),
	}

	// readDeadline := time.Now().Add(time.Duration(connection.Conf.ReadTimeout))
	// writeDeadline := time.Now().Add(time.Duration(connection.Conf.WriteTimeout))

	// if err := connection.conn.SetReadDeadline(readDeadline); err != nil {
	// 	return nil, err
	// }

	// err := connection.conn.SetWriteDeadline(writeDeadline)
	return connection, nil
}

func (c *Conn) Read() ([]byte, error) {
	byteLength, err := c.conn.Read(c.chunk)
	if err != nil {
		log.Println("read failed, err:", err)
		return nil, err
	}

	if int64(byteLength) < c.Conf.MinByteSize {
		return nil, ErrByteLengthTooSmall
	}

	msg := make([]byte, byteLength)
	copy(msg, c.chunk)
	return msg, nil
}

func (c *Conn) Write(msg []byte) error {
	byteLength, err := c.conn.Write(msg)
	if err != nil {
		log.Println("write failed, err:", err)
		return err
	}

	log.Println("write success, byteLength:", byteLength)
	return nil
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
