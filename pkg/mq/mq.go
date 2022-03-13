package mq

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
)

const _network = "tcp"

const (
	_publisher = iota + 1
	_subscriber
)

type (
	ReaderConfig struct {
		URL         string `json:"url"`
		Topic       string `json:"topic"`
		ReadTimeout int64  `json:"read_timeout"`
		MaxByteSize int64  `json:"max_byte_size"`
		MinByteSize int64  `json:"min_byte_size"`
	}

	Reader struct {
		conn net.Conn
	}
)

func (rc ReaderConfig) toConnectionConfig() ConnectionConfig {
	return ConnectionConfig{
		Topic:       rc.Topic,
		ReadTimeout: rc.ReadTimeout,
		MaxByteSize: rc.MaxByteSize,
		MinByteSize: rc.MinByteSize,
		Type:        _subscriber,
	}
}

func NewReader(conf ReaderConfig) (*Reader, error) {
	conn, err := connect(conf.URL)
	if err != nil {
		return nil, err
	}

	err = handshake(conn, conf.toConnectionConfig())
	return &Reader{conn}, err
}

// Read TODO: this method needs to be updated now it is not useful
func (r *Reader) Read() {
	chunk := make([]byte, 2048)
	for {
		length, err := r.conn.Read(chunk)
		if err != nil {
			log.Printf("msg could not read: %v\n", err)
			if err := r.conn.Close(); err != nil {
				log.Printf("could not close connection: %v\n", err)
			}
			break
		}

		msg := make([]byte, length)
		copy(msg, chunk)

		log.Printf("length: %d, msg: %s\n", length, string(msg))
	}

	if err := r.conn.Close(); err != nil {
		log.Printf("could not close connection: %v\n", err)
	}
}

type (
	WriterConfig struct {
		URL          string `json:"url"`
		Topic        string `json:"topic"`
		WriteTimeout int64  `json:"write_timeout"`
		MaxByteSize  int64  `json:"max_byte_size"`
		MinByteSize  int64  `json:"min_byte_size"`
	}

	Writer struct {
		conn net.Conn
	}
)

func (wc WriterConfig) toConnectionConfig() ConnectionConfig {
	return ConnectionConfig{
		Topic:        wc.Topic,
		WriteTimeout: wc.WriteTimeout,
		MaxByteSize:  wc.MaxByteSize,
		MinByteSize:  wc.MinByteSize,
		Type:         _publisher,
	}
}

func NewWriter(conf WriterConfig) (*Writer, error) {
	conn, err := connect(conf.URL)
	if err != nil {
		return nil, err
	}

	err = handshake(conn, conf.toConnectionConfig())
	return &Writer{conn}, err
}

// Write TODO: This method needs to be updated now it is not useful
func (w *Writer) Write() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		var msgContent string
		if scanner.Scan() {
			msgContent = scanner.Text()
		}

		if msgContent == "q" {
			break
		}

		if _, err := w.conn.Write([]byte(msgContent)); err != nil {
			log.Println("write failed, err:", err)
		}
	}

	if err := w.conn.Close(); err != nil {
		log.Println("could not close connection:", err)
	}
}

type ConnectionConfig struct {
	Topic        string `json:"topic"`
	ReadTimeout  int64  `json:"read_timeout"`
	WriteTimeout int64  `json:"write_timeout"`
	MaxByteSize  int64  `json:"max_byte_size"`
	MinByteSize  int64  `json:"min_byte_size"`
	Type         uint8  `json:"type"`
}

func connect(url string) (net.Conn, error) {
	return net.Dial(_network, url)
}

func handshake(conn net.Conn, conf ConnectionConfig) error {
	// read health message
	chunk := make([]byte, 1024)
	if _, err := conn.Read(chunk); err != nil {
		return err
	}
	log.Println(string(chunk))

	// write configuration message
	confBytes, _ := json.Marshal(conf)
	_, err := conn.Write(confBytes)
	return err
}
