package conn

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
)

const _network = "tcp"

var (
	healthyHelloMsg = []byte(`{"healthy": true}`)
	chunk           = make([]byte, 1024)
)

type Server struct {
	listener    net.Listener
	connections []Conn

	OnConnection            func(c Conn)
	ConnectionWriteDeadline time.Duration
	ConnectionReadDeadline  time.Duration
}

// Start a server on the given address.
func (s *Server) Start(addr string) error {
	listener, err := net.Listen(_network, addr)
	if err != nil {
		return err
	}

	s.listener = listener
	return nil
}

// Listen for incoming connections.
func (s *Server) Listen() {
	for {
		stdconn, err := s.listener.Accept()
		if err != nil {
			log.Printf("connection dropped, err: %v\n", err)
		}

		go s.acceptConnection(stdconn)
	}
}

func (s *Server) acceptConnection(stdconn net.Conn) {
	writeDeadline := time.Now().Add(s.ConnectionWriteDeadline)
	if err := stdconn.SetWriteDeadline(writeDeadline); err != nil {
		log.Println("write deadline exceeded")
	}

	readDeadline := time.Now().Add(s.ConnectionReadDeadline)
	if err := stdconn.SetReadDeadline(readDeadline); err != nil {
		log.Println("read deadline exceeded")
	}

	if _, err := stdconn.Write(healthyHelloMsg); err != nil {
		log.Printf("write failed, err: %v\n", err)
	}
	log.Println("health message sent")

	byteLength, err := stdconn.Read(chunk)
	if err != nil {
		log.Printf("read failed, msg: %s, err: %v\n", string(chunk), err)
	}

	msg := make([]byte, byteLength)
	copy(msg, chunk)

	var conf Config
	if err := json.Unmarshal(msg, &conf); err != nil {
		log.Printf("unmarshal failed, msg: %s, err: %v\n", string(msg), err)
		if _, err := stdconn.Write([]byte(`{"error": "incorrect config message format"}`)); err != nil {
			log.Printf("write failed, err: %v\n", err)
			return
		}
	}

	conn := Conn{
		ID:   uuid.NewString(),
		conn: stdconn,
		At:   time.Now(),
		Conf: conf,
	}
	s.connections = append(s.connections, conn)

	s.OnConnection(conn)

}

func (s *Server) Close() error {
	return s.listener.Close()
}

type (
	Config struct {
		Topic        string `json:"topic"`
		ReadTimeout  int64  `json:"read_timeout"`
		WriteTimeout int64  `json:"write_timeout"`
		MaxByteSize  int64  `json:"max_byte_size"`
		MinByteSize  int64  `json:"min_byte_size"`

		// Encryption
	}

	Conn struct {
		conn net.Conn

		ID   string
		At   time.Time
		Conf Config
	}
)

func (c *Conn) Close() error {
	return c.conn.Close()
}
