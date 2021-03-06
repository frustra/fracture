package edge

import (
	"errors"
	"log"
	"math"
	"math/rand"
	"net"

	"github.com/frustra/fracture/edge/protocol"
	"github.com/frustra/fracture/network"
	"github.com/frustra/fracture/protobuf"
)

type Server struct {
	Addr    string
	Offset  int64
	Size    int
	Cluster *network.Cluster

	keyPair *protocol.KeyPair

	EntityServers     map[*network.InternalConnection]int
	Clients           map[*GameConnection]bool
	PlayerConnections map[string]chan *protocol.Packet
}

func (s *Server) HandleMessage(message interface{}, conn *network.InternalConnection) {
	switch msg := message.(type) {
	case *protobuf.ChunkResponse:
		s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.MapChunkBulkID, int16(1), int32(len(msg.Data)), true, msg.Data, int32(msg.X), int32(msg.Z), uint16(0xFFFF), uint16(0))
	case *protobuf.ChatMessage:
		s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.ChatMessageID, msg.Message)
	case *protobuf.BlockUpdate:
		log.Printf("Sending block update to: %s", msg.Uuid)
		s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.BlockChangeID, int32(msg.X), uint8(msg.Y), int32(msg.Z), protocol.Varint{uint64(msg.BlockId)}, uint8(msg.BlockMetadata))
	case *protobuf.PlayerAction:
		switch msg.Action {
		case protobuf.PlayerAction_JOIN:
			s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.PlayerListItemID, msg.Player.Username, true, int16(0))
			s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.ChatMessageID, protocol.CreateJsonMessage(msg.Player.Username+" joined the game", "yellow"))
			if msg.Uuid != msg.Player.Uuid {
				meta := protocol.NewMetadata(
					protocol.AbsorptionHeartsID, float32(0),
					protocol.OnFireID, false,
					protocol.UnknownBitFieldID, byte(0),
					protocol.AirID, uint16(0x012c),
					protocol.ScoreID, uint32(0),
					protocol.HealthID, float32(20),
					protocol.PotionColorID, int32(0),
					protocol.AmbientPotionID, byte(0),
					protocol.ArrowCountID, byte(0),
				)

				s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.SpawnPlayerID,
					protocol.Varint{uint64(msg.Player.EntityId)},
					msg.Player.Uuid,
					msg.Player.Username,
					protocol.Varint{0},
					int32(msg.Player.X*32),
					int32(msg.Player.FeetY*32),
					int32(msg.Player.Z*32),
					byte(msg.Player.Yaw*256/2/math.Pi),
					byte(msg.Player.Pitch*256/2/math.Pi),
					int16(0),
					meta,
				)
				s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.EntityMetadataID,
					int32(msg.Player.EntityId),
					meta,
				)
				s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.EntityPropertiesID,
					int32(msg.Player.EntityId),
					int32(2),
					"generic.maxHealth", float64(20), int16(0),
					"generic.movementSpeed", float64(0.1), int16(0),
				)
				s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.EntityTeleportID,
					int32(msg.Player.EntityId),
					int32(msg.Player.X*32),
					int32(msg.Player.FeetY*32),
					int32(msg.Player.Z*32),
					byte(msg.Player.Yaw*256/360),
					byte(msg.Player.Pitch*256/360),
				)
			}
		case protobuf.PlayerAction_MOVE_RELATIVE:
			if msg.Flags == 1 {
				s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.EntityRelativeMoveID,
					int32(msg.Player.EntityId),
					byte(msg.Player.X*32),
					byte(msg.Player.FeetY*32),
					byte(msg.Player.Z*32),
				)
			} else if msg.Flags == 2 {
				s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.EntityLookID,
					int32(msg.Player.EntityId),
					byte(msg.Player.Yaw*256/360),
					byte(msg.Player.Pitch*256/360),
				)
			} else if msg.Flags == 3 {
				s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.EntityLookAndMoveID,
					int32(msg.Player.EntityId),
					byte(msg.Player.X*32),
					byte(msg.Player.FeetY*32),
					byte(msg.Player.Z*32),
					byte(msg.Player.Yaw*256/360),
					byte(msg.Player.Pitch*256/360),
				)
			}
			if msg.Flags&2 == 2 {
				s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.EntityHeadLookID, int32(msg.Player.EntityId), byte(msg.Player.Yaw*256/360))
			}
		case protobuf.PlayerAction_MOVE_ABSOLUTE:
			s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.EntityTeleportID,
				int32(msg.Player.EntityId),
				int32(msg.Player.X*32),
				int32(msg.Player.FeetY*32),
				int32(msg.Player.Z*32),
				byte(msg.Player.Yaw*256/360),
				byte(msg.Player.Pitch*256/360),
			)
			s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.EntityHeadLookID, int32(msg.Player.EntityId), byte(msg.Player.Yaw*256/360))
		case protobuf.PlayerAction_LEAVE:
			s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.DestroyEntitiesID, byte(1), int32(msg.Player.EntityId))
			s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.PlayerListItemID, msg.Player.Username, false, int16(0))
			s.PlayerConnections[msg.Uuid] <- protocol.CreatePacket(protocol.ChatMessageID, protocol.CreateJsonMessage(msg.Player.Username+" left the game", "yellow"))
		}
	}
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
	s.PlayerConnections = make(map[string]chan *protocol.Packet)

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
		go client.HandleAuth()
	}
}

func (s *Server) NodePort() int {
	return 0
}

func (s *Server) GetMinecraftStatus() protocol.StatusResponse {
	statusMessage := protocol.JsonMessage{"Fracture Distributed Server", "green", false, false, false, false, false}
	return protocol.CreateStatusResponse("1.7.10", protocol.Version, 0, s.Size, statusMessage)
}

func (s *Server) DrainPlayerConnections(cc *GameConnection) {
	for {
		msg := <-s.PlayerConnections[cc.Player.Uuid]
		if msg == nil {
			return
		}
		msg.Write(cc.ConnEncrypted)
	}
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
		entityServers := s.Cluster.EntityNodes
		serverRange := make([]string, len(entityServers))
		i := 0
		for _, node := range entityServers {
			serverRange[i] = node.Meta.Addr
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
