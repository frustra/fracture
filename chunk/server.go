package chunk

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

	OffsetX, OffsetZ int64

	Storage *Chunk
}

func (s *Server) Serve() error {
	log.Printf("Chunk server (%d, %d) loading on %s\n", s.OffsetX, s.OffsetZ, s.Addr)
	s.Storage = NewChunk(s.OffsetX, s.OffsetZ)
	return network.ServeInternal(s.Addr, s)
}

func (s *Server) HandleMessage(message interface{}, conn *network.InternalConnection) {
	switch req := message.(type) {
	case *protobuf.ChunkRequest:
		x, z := req.GetX(), req.GetZ()

		res := &protobuf.ChunkResponse{
			X:    x,
			Z:    z,
			Data: s.Storage.MarshallCompressed(),
		}

		conn.SendMessage(res)
	}
}

func (s *Server) NodeType() string {
	return "chunk"
}

func (s *Server) NodePort() int {
	_, metaPortString, _ := net.SplitHostPort(s.Addr)
	port, _ := strconv.Atoi(metaPortString)
	return port
}
