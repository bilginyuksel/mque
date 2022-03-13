package conn

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

const _network = "tcp"

var (
	healthyHelloMsg = []byte(`{"status": "healthy"}`)
	chunk           = make([]byte, 1024)
)

type Server struct {
	listener    net.Listener
	Connections map[string]*Conn

	OnConnection            func(c *Conn)
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
	// writeDeadline := time.Now().Add(s.ConnectionWriteDeadline)
	// if err := stdconn.SetWriteDeadline(writeDeadline); err != nil {
	// 	log.Println("write deadline exceeded")
	// }

	// readDeadline := time.Now().Add(s.ConnectionReadDeadline)
	// if err := stdconn.SetReadDeadline(readDeadline); err != nil {
	// 	log.Println("read deadline exceeded")
	// }

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

	conn, err := New(stdconn, conf)
	if err != nil {
		log.Printf("new connection failed, err: %v\n", err)
		return
	}
	s.Connections[conn.ID] = conn

	s.OnConnection(conn)
}

func (s *Server) Close() error {
	return s.listener.Close()
}
