package edge

import (
	"fmt"
	"net"
	"os"

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

func (s *Server) Serve() {
	addr, err := net.ResolveTCPAddr("tcp4", s.Addr)
	assertNoErr(err)

	listener, err := net.ListenTCP("tcp", addr)
	assertNoErr(err)

	tmpkey, err := protocol.GenerateKeyPair(1024)
	assertNoErr(err)
	s.keyPair = tmpkey

	fmt.Printf("Listening for TCP on %s\n", s.Addr)
	defer listener.Close()

	s.Clients = make(map[*GameConnection]bool)

	for {
		conn, err := listener.AcceptTCP()
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

func assertNoErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err.Error())
		os.Exit(1)
	}
}
