package edge

import (
	"log"
	"net"

	"github.com/frustra/fracture/network"
)

type Server struct {
	Addr    string
	Cluster *network.Cluster

	OffsetX, OffsetY int
	Size             int // Inclusive diameter.
}

func (s *Server) Serve() error {
	log.Printf("Chunk server loading on %s\n", s.Addr)
	return network.ServeInternal(s.Addr, s)
}

func (s *Server) NodeType() string {
	return "chunk"
}

func (s *Server) NodePort() int {
	return 1234
}
