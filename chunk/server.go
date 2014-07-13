package chunk

import (
	"log"
	"net"
	"strconv"

	"github.com/frustra/fracture/network"
	"github.com/frustra/fracture/protobuf"
)

const (
	ChunkWidthPerNode = 8 // Side length of a square.
)

type Server struct {
	Addr    string
	Cluster *network.Cluster

	OffsetX, OffsetZ int64

	Storage [ChunkWidthPerNode][ChunkWidthPerNode]*Chunk // [z][x]
}

func (s *Server) Serve() error {
	log.Printf("Chunk server (%d, %d) loading on %s\n", s.OffsetX, s.OffsetZ, s.Addr)

	blockType := byte(((s.OffsetX+8)/8+(s.OffsetZ+8)/4)%4 + 1)

	for z := int64(0); z < ChunkWidthPerNode; z++ {
		for x := int64(0); x < ChunkWidthPerNode; x++ {
			c := NewChunk(s.OffsetX+x, s.OffsetZ+z)
			c.Generate(blockType)
			s.Storage[z][x] = c
		}
	}
	return network.ServeInternal(s.Addr, s)
}

func (s *Server) HandleMessage(message interface{}, conn *network.InternalConnection) {
	switch req := message.(type) {
	case *protobuf.ChunkRequest:
		x, z := req.X-s.OffsetX, req.Z-s.OffsetZ

		res := &protobuf.ChunkResponse{
			X:    req.X,
			Z:    req.Z,
			Uuid: req.Uuid,
		}

		if x < 0 || z < 0 || x >= ChunkWidthPerNode || z >= ChunkWidthPerNode {
			res.Data = make([]byte, 0)
		} else {
			res.Data = s.Storage[z][x].MarshallCompressed()
		}

		conn.SendMessage(res)
	case *protobuf.BlockUpdate:
		cx, cz := WorldCoordsToChunk(req.X, req.Z)
		cx -= s.OffsetX
		cz -= s.OffsetZ

		if cx < 0 || cz < 0 || cx >= ChunkWidthPerNode || cz >= ChunkWidthPerNode {
			log.Printf("Received block update for someone else's chunk: %#v", req)
		} else {
			chunk := s.Storage[cz][cx]
			x, z := WorldCoordsToChunkInternal(req.X, req.Z)
			y := req.Y

			if req.Destroy {
				chunk.Set(int(x), int(y), int(z), 0)
			} else {
				panic("unimplemented")
			}
		}
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
