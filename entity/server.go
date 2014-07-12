package entity

import (
	"log"
	"net"
	"strconv"

	"github.com/frustra/fracture/network"
	"github.com/frustra/fracture/protobuf"
)

type Server struct {
	Addr    string
	Cluster *network.Cluster

	Size  int
	XSort []*protobuf.Player
	YSort []*protobuf.Player
}

func (s *Server) Serve() error {
	log.Printf("Entity server loading on %s\n", s.Addr)
	return network.ServeInternal(s.Addr, s)
}

func (s *Server) HandleMessage(message interface{}, conn *network.InternalConnection) {
	log.Print("Handler invoked: ", message)
}

func (s *Server) NodeType() string {
	return "entity"
}

func (s *Server) NodePort() int {
	_, metaPortString, _ := net.SplitHostPort(s.Addr)
	port, _ := strconv.Atoi(metaPortString)
	return port
}
