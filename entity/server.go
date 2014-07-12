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

	Size    int
	Players map[string]*PlayerEntity
}

type PlayerEntity struct {
	protobuf.Player
	Conn *network.InternalConnection
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
					X:        0,
					HeadY:    128,
					FeetY:    128 - 1.62,
					Z:        0,
				},
				Conn: conn,
			}

			log.Printf("Player joined: %s", player.Username)

			for _, p := range s.Players {
				p.Conn.SendMessage(&protobuf.PlayerAction{
					Player: &s.Players[player.Uuid].Player,
					Action: protobuf.PlayerAction_JOIN,
				})
			}
		case protobuf.PlayerAction_MOVE:
			tmp := s.Players[player.Uuid]
			tmp.X = player.X
			tmp.HeadY = player.HeadY
			tmp.FeetY = player.FeetY
			tmp.Z = player.Z
			tmp.Pitch = player.Pitch
			tmp.Yaw = player.Yaw

			for uuid, p := range s.Players {
				if uuid != player.Uuid {
					p.Conn.SendMessage(&protobuf.PlayerAction{
						Player: &protobuf.Player{
							Uuid:  player.Uuid,
							X:     player.X,
							HeadY: player.HeadY,
							FeetY: player.FeetY,
							Z:     player.Z,
							Pitch: player.Pitch,
							Yaw:   player.Yaw,
						},
						Action: protobuf.PlayerAction_MOVE,
					})
				}
			}
		case protobuf.PlayerAction_LEAVE:
			delete(s.Players, player.Uuid)

			log.Printf("Player left: %s", player.Username)

			for _, p := range s.Players {
				p.Conn.SendMessage(&protobuf.PlayerAction{
					Player: &protobuf.Player{Uuid: player.Uuid},
					Action: protobuf.PlayerAction_LEAVE,
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
