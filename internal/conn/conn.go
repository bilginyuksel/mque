package con

import (
	"log"
	"net"
	"time"

	"github.com/google/uuid"
)

const _network = "tcp"

type Server struct {
	listener     net.Listener
	connections  []Conn
	OnConnection func(c Conn)
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
	stdconn, err := s.listener.Accept()
	if err != nil {
		log.Printf("connection could not be accepted, err: %v\n", err)
	}

	conn := Conn{
		id:   uuid.NewString(),
		conn: stdconn,
		at:   time.Now(),
	}
	s.connections = append(s.connections, conn)

	s.OnConnection(conn)
}

func (s *Server) Close() error {
	return s.listener.Close()
}

type Conn struct {
	id   string
	conn net.Conn
	at   time.Time
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
