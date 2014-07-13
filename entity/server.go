package entity

import (
	"log"
	"math"
	"net"
	"strconv"

	"github.com/frustra/fracture/edge/protocol"
	"github.com/frustra/fracture/network"
	"github.com/frustra/fracture/protobuf"
)

type Server struct {
	Addr    string
	Cluster *network.Cluster

	Size    int
	Players map[string]*PlayerEntity
}

type PlayerEntity struct {
	protobuf.Player
	LastX float64
	LastY float64
	LastZ float64
	Conn  *network.InternalConnection
}

func (s *Server) Serve() error {
	log.Printf("Entity server loading on %s\n", s.Addr)

	s.Players = make(map[string]*PlayerEntity)

	return network.ServeInternal(s.Addr, s)
}

func (s *Server) HandleMessage(message interface{}, conn *network.InternalConnection) {
	switch req := message.(type) {
	case *protobuf.ChatMessage:
		player := s.Players[req.Uuid]
		message := protocol.CreateJsonMessage("<"+player.Username+"> "+req.Message, "")
		log.Printf("Chat: %s", message.Message)
		for uuid, p := range s.Players {
			p.Conn.SendMessage(&protobuf.ChatMessage{
				Message: message.String(),
				Uuid:    uuid,
			})
		}
	case *protobuf.BlockUpdate:
		for uuid, p := range s.Players {
			p.Conn.SendMessage(&protobuf.BlockUpdate{
				X:             req.X,
				Y:             req.Y,
				Z:             req.Z,
				BlockId:       req.BlockId,
				BlockMetadata: req.BlockMetadata,
				Uuid:          uuid,
			})
		}
	case *protobuf.PlayerAction:
		player, action := req.GetPlayer(), req.GetAction()
		switch action {
		case protobuf.PlayerAction_JOIN:
			s.Players[player.Uuid] = &PlayerEntity{
				Player: protobuf.Player{
					Uuid:     player.Uuid,
					Username: player.Username,
					EntityId: player.EntityId,
					X:        0,
					HeadY:    105,
					FeetY:    105 - 1.62,
					Z:        0,
				},
				LastX: 0,
				LastY: 105 - 1.62,
				LastZ: 0,
				Conn:  conn,
			}

			log.Printf("Player joined (%d): %s", player.EntityId, player.Username)

			for uuid, p := range s.Players {
				p.Conn.SendMessage(&protobuf.PlayerAction{
					Player: &s.Players[player.Uuid].Player,
					Action: protobuf.PlayerAction_JOIN,
					Uuid:   uuid,
				})
				if uuid != player.Uuid {
					conn.SendMessage(&protobuf.PlayerAction{
						Player: &p.Player,
						Action: protobuf.PlayerAction_JOIN,
						Uuid:   player.Uuid,
					})
				}
			}
		case protobuf.PlayerAction_MOVE_ABSOLUTE:
			responseType := protobuf.PlayerAction_MOVE_RELATIVE

			tmp := s.Players[player.Uuid]
			if req.Flags&1 == 1 {
				if math.Abs(player.X-tmp.LastX) > 3 ||
					math.Abs(player.FeetY-tmp.LastY) > 3 ||
					math.Abs(player.Z-tmp.LastZ) > 3 {
					responseType = protobuf.PlayerAction_MOVE_ABSOLUTE
					tmp.LastX = player.X
					tmp.LastY = player.FeetY
					tmp.LastZ = player.Z
				}
			}

			for uuid, p := range s.Players {
				if uuid != player.Uuid {
					sendPlayer := protobuf.Player{
						Uuid:     player.Uuid,
						EntityId: tmp.EntityId,
					}
					if responseType == protobuf.PlayerAction_MOVE_ABSOLUTE {
						sendPlayer.X = player.X
						sendPlayer.HeadY = player.HeadY
						sendPlayer.FeetY = player.FeetY
						sendPlayer.Z = player.Z
					} else if req.Flags&1 == 1 {
						sendPlayer.X = player.X - tmp.X
						sendPlayer.FeetY = player.FeetY - tmp.FeetY
						sendPlayer.Z = player.Z - tmp.Z
					}
					if req.Flags&2 == 2 {
						sendPlayer.Pitch = player.Pitch
						sendPlayer.Yaw = player.Yaw
					}
					p.Conn.SendMessage(&protobuf.PlayerAction{
						Player: &sendPlayer,
						Action: responseType,
						Uuid:   uuid,
						Flags:  req.Flags,
					})
				}
			}

			if req.Flags&1 == 1 {
				tmp.X = player.X
				tmp.HeadY = player.HeadY
				tmp.FeetY = player.FeetY
				tmp.Z = player.Z
			}
			if req.Flags&2 == 2 {
				tmp.Pitch = player.Pitch
				tmp.Yaw = player.Yaw
			}
		case protobuf.PlayerAction_LEAVE:
			tmp := s.Players[player.Uuid]
			delete(s.Players, player.Uuid)

			log.Printf("Player left: %s", tmp.Username)

			for uuid, p := range s.Players {
				sendPlayer := protobuf.Player{
					Uuid:     player.Uuid,
					Username: tmp.Username,
					EntityId: tmp.EntityId,
				}
				p.Conn.SendMessage(&protobuf.PlayerAction{
					Player: &sendPlayer,
					Action: protobuf.PlayerAction_LEAVE,
					Uuid:   uuid,
				})
			}
		}
	}
}

func (s *Server) NodeType() string {
	return "entity"
}

func (s *Server) NodePort() int {
	_, metaPortString, _ := net.SplitHostPort(s.Addr)
	port, _ := strconv.Atoi(metaPortString)
	return port
}
