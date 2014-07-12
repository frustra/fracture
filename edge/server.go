package edge

import (
	"errors"
	"log"
	"math"
	"math/rand"
	"net"

	"github.com/frustra/fracture/chunk"
	"github.com/frustra/fracture/edge/protocol"
	"github.com/frustra/fracture/network"
	"github.com/frustra/fracture/protobuf"
)

type ChunkCoord struct {
	X, Z int64
}

type Server struct {
	Addr    string
	Size    int
	Cluster *network.Cluster

	keyPair *protocol.KeyPair

	EntityServers map[*network.InternalConnection]int
	ChunkServers  map[ChunkCoord]*network.InternalConnection
	Clients       map[*GameConnection]bool
}

func (s *Server) HandleMessage(message interface{}, conn *network.InternalConnection) {
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
	s.EntityServers = make(map[*network.InternalConnection]int)
	s.ChunkServers = make(map[ChunkCoord]*network.InternalConnection)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		client := &GameConnection{
			Server: s,
			Conn:   conn,
			Player: &protobuf.Player{},
		}
		s.Clients[client] = true
		go client.HandleConnection()
	}
}

func (s *Server) NodeType() string {
	return "edge"
}

func (s *Server) NodePort() int {
	return 0
}

func (s *Server) FindEntityServer(player *protobuf.Player) (*network.InternalConnection, error) {
	var closestDist float64 = -1
	var closestServer *network.InternalConnection
	for client, _ := range s.Clients {
		if client.EntityServer != nil {
			dist := math.Abs(client.Player.X - player.X)
			if closestDist < 0 || closestDist > dist {
				closestDist = dist
				closestServer = client.EntityServer
			}
		}
	}
	if closestServer != nil {
		return closestServer, nil
	} else {
		entityServers := s.Cluster.MetaLookup["entity"]
		serverRange := make([]string, len(entityServers))
		i := 0
		for _, meta := range entityServers {
			serverRange[i] = meta.GetAddr()
			i++
		}
		if i > 0 {
			addr := serverRange[rand.Intn(i)]
			return network.ConnectInternal(addr, s)
		} else {
			return nil, errors.New("No entity servers available!")
		}
	}
}

func (s *Server) FindChunkServer(x, z int64) (*network.InternalConnection, error) {
	x, z = chunk.ChunkCoordsToNode(x, z)
	coord := ChunkCoord{x, z}

	if conn, exists := s.ChunkServers[coord]; exists {
		return conn, nil
	}

	chunkServers := s.Cluster.MetaLookup["chunk"]
	for _, meta := range chunkServers {
		if *meta.X == x && *meta.Z == z {
			conn, err := network.ConnectInternal(meta.GetAddr(), s)
			if err != nil {
				return nil, err
			}
			s.ChunkCoord[coord] = conn
			return conn, nil
		}
	}
	return nil, errors.New("No chunk server for this area!")
}
