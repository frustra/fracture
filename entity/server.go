package entity

import (
	"log"

	"github.com/frustra/fracture/network"
)

type Server struct {
	Addr    string
	Cluster *network.Cluster

	Size  int
	XSort []*Player
	YSort []*Player
}

func (s *Server) Serve() error {
	log.Printf("Entity server loading on %s\n", s.Addr)
	return network.ServeInternal(s.Addr, s)
}

func (s *Server) HandleMessage(message interface{}) {
	log.Print("Handler invoked: ", message)
}

func (s *Server) NodeType() string {
	return "entity"
}

func (s *Server) NodePort() int {
	return 1234
}
