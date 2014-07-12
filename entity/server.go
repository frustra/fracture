package entity

import (
	"log"
	"math"
	"net"
	"strconv"

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

			log.Printf("Player joined: %s", player.Username)

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
			if math.Abs(player.X-tmp.LastX) > 3 ||
				math.Abs(player.FeetY-tmp.LastY) > 3 ||
				math.Abs(player.Z-tmp.LastZ) > 3 {
				responseType = protobuf.PlayerAction_MOVE_ABSOLUTE
				tmp.LastX = player.X
				tmp.LastY = player.FeetY
				tmp.LastZ = player.Z
			}

			for uuid, p := range s.Players {
				if uuid != player.Uuid {
					sendPlayer := protobuf.Player{
						Uuid:     player.Uuid,
						EntityId: player.EntityId,
						Pitch:    player.Pitch,
						Yaw:      player.Yaw,
					}
					if responseType == protobuf.PlayerAction_MOVE_ABSOLUTE {
						sendPlayer.X = player.X
						sendPlayer.HeadY = player.HeadY
						sendPlayer.FeetY = player.FeetY
						sendPlayer.Z = player.Z
					} else {
						sendPlayer.X = player.X - tmp.X
						sendPlayer.FeetY = player.FeetY - tmp.FeetY
						sendPlayer.Z = player.Z - tmp.Z
					}
					p.Conn.SendMessage(&protobuf.PlayerAction{
						Player: &sendPlayer,
						Action: responseType,
						Uuid:   uuid,
					})
				}
			}

			tmp.X = player.X
			tmp.HeadY = player.HeadY
			tmp.FeetY = player.FeetY
			tmp.Z = player.Z
			tmp.Pitch = player.Pitch
			tmp.Yaw = player.Yaw
		case protobuf.PlayerAction_LEAVE:
			delete(s.Players, player.Uuid)

			log.Printf("Player left: %s", player.Username)

			for uuid, p := range s.Players {
				p.Conn.SendMessage(&protobuf.PlayerAction{
					Player: &protobuf.Player{Uuid: player.Uuid, Username: player.Username, EntityId: player.EntityId},
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
