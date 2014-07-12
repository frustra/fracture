package edge

import (
	"log"
	"math/rand"
	"net"

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

	for {
		conn, err := listener.Accept()
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
func (s *Server) FindEntityServer(x float64) string {
	var minPlayer *GameConnection // One player right of x
	var maxPlayer *GameConnection // One player left of x
	for client, _ := range s.Clients {
		if client.Connected {
			if client.X >= x {
				if minPlayer == nil || minPlayer.X < client.X {
					minPlayer = client
				}
			}
			if client.X <= x {
				if maxPlayer == nil || maxPlayer.X > client.X {
					maxPlayer = client
				}
			}
		}
	}
	entityServers := s.Cluster.MetaLookup["entity"]
	serverRange := make([]string, len(entityServers))
	i := 0
	for _, meta := range entityServers {
		if meta.GetX() >= maxPlayer.ServerId && meta.GetX() <= minPlayer.ServerId {
			serverRange[i] = meta.GetAddr()
			i++
		}
	}
	return serverRange[rand.Intn(i)]
}
