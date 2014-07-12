package chunk

import (
	"log"

	"github.com/frustra/fracture/network"
	"github.com/frustra/fracture/protobuf"
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

func (s *Server) HandleMessage(message interface{}, conn *network.InternalConnection) {
	switch req := message.(type) {
	case *protobuf.ChunkRequest:
		x, z := req.GetX(), req.GetZ()

		res := &protobuf.ChunkResponse{
			X:    &x,
			Z:    &z,
			Data: make([]byte, 0),
		}

		conn.SendMessage(res)
	}
}

func (s *Server) NodeType() string {
	return "chunk"
}

func (s *Server) NodePort() int {
	return 1234
}
