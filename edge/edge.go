package edge

import (
	"github.com/frustra/fracture/network"
)

type Server struct {
	Addr       string
	MaxPlayers int
	Cluster    *network.Cluster
}

func (s *Server) Serve() {
	var ch chan bool
	ch <- true
}

func (s *Server) NodeType() string {
	return "edge"
}

func (s *Server) NodePort() int {
	return 1234
}
