package edge

import (
	"log"
	"net"

	"github.com/frustra/fracture/edge/protocol"
	"github.com/frustra/fracture/network"
)

type Server struct {
	Addr       string
	MaxPlayers int
	Cluster    *network.Cluster

	keyPair *protocol.KeyPair

	Clients map[*GameConnection]bool
}

func (s *Server) Serve() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	tmpkey, err := protocol.GenerateKeyPair(1024)
	if err != nil {
		return err
	}
	s.keyPair = tmpkey

	log.Printf("Game connection listening on %s\n", s.Addr)
	defer listener.Close()

	s.Clients = make(map[*GameConnection]bool)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		client := &GameConnection{s, conn, nil, true, ""}
		s.Clients[client] = true
		go client.HandleConnection()
	}
}

func (s *Server) NodeType() string {
	return "edge"
}

func (s *Server) NodePort() int {
	return 1234
}
