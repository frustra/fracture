package chunk

import (
	"log"
	"net"
	"strconv"

	"github.com/frustra/fracture/network"
	"github.com/frustra/fracture/protobuf"
	"github.com/frustra/fracture/world"
)

type Server struct {
	Addr    string
	Cluster *network.Cluster

	OffsetX, OffsetZ int64
	Listeners        map[*network.InternalConnection]bool

	storage [world.ChunkWidthPerNode][world.ChunkWidthPerNode]*world.Chunk // [z][x]

	blockUpdateChannel chan *protobuf.BlockUpdate
}

func (s *Server) Serve() error {
	log.Printf("Chunk server (%d, %d) loading on %s\n", s.OffsetX, s.OffsetZ, s.Addr)

	s.Listeners = make(map[*network.InternalConnection]bool)

	blockType := byte(((s.OffsetX+8)/8+(s.OffsetZ+8)/4)%4 + 1)

	for z := int64(0); z < world.ChunkWidthPerNode; z++ {
		for x := int64(0); x < world.ChunkWidthPerNode; x++ {
			c := world.NewChunk(s.OffsetX+x, s.OffsetZ+z)
			c.Generate(blockType)
			s.storage[z][x] = c
		}
	}

	s.blockUpdateChannel = make(chan *protobuf.BlockUpdate)
	go s.Loop()

	return network.ServeInternal(s.Addr, s)
}

func (s *Server) Loop() {
	for {
		select {
		case update := <-s.blockUpdateChannel:
			chunk := s.GetChunk(world.WorldCoordsToChunk(update.X, update.Z))

			if chunk == nil {
				log.Printf("Received block update for someone else's chunk: %#v", update)
				continue
			}

			x, z := world.WorldCoordsToChunkInternal(update.X, update.Z)
			y := update.Y

			chunk.Set(x, int64(y), z, byte(update.BlockId))
			chunk.SetMetadata(x, int64(y), z, byte(update.BlockMetadata))
			chunk.CalculateSkyLightingForColumn(x, z)

			for conn, _ := range s.Listeners {
				conn.SendMessage(update)
			}
		}
	}
}

func (s *Server) LocalChunkCoords(x, z int64) (int64, int64) {
	return x - s.OffsetX, z - s.OffsetZ
}

func (s *Server) ContainsLocalChunk(x, z int64) bool {
	return x >= 0 && z >= 0 && x < world.ChunkWidthPerNode && z < world.ChunkWidthPerNode
}

func (s *Server) GetChunk(x, z int64) *world.Chunk {
	x, z = s.LocalChunkCoords(x, z)

	if s.ContainsLocalChunk(x, z) {
		return s.storage[z][x]
	}

	return nil
}

func (s *Server) HandleMessage(message interface{}, conn *network.InternalConnection) {
	switch req := message.(type) {
	case *protobuf.Subscription:
		if req.Subscribe {
			log.Printf("Got subscription from (%d, %d): %s", s.OffsetX, s.OffsetZ, conn)
			s.Listeners[conn] = true
		} else {
			delete(s.Listeners, conn)
		}
	case *protobuf.ChunkRequest:
		chunk := s.GetChunk(req.X, req.Z)

		res := &protobuf.ChunkResponse{
			X:    req.X,
			Z:    req.Z,
			Uuid: req.Uuid,
		}

		if chunk != nil {
			res.Data = chunk.MarshallCompressed()
		} else {
			res.Data = make([]byte, 0)
		}

		conn.SendMessage(res)
	case *protobuf.BlockUpdate:
		s.blockUpdateChannel <- req
	}
}

func (s *Server) NodePort() int {
	_, metaPortString, _ := net.SplitHostPort(s.Addr)
	port, _ := strconv.Atoi(metaPortString)
	return port
}
