package edge

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/frustra/fracture/chunk"
	"github.com/frustra/fracture/edge/protocol"
	"github.com/frustra/fracture/network"
	"github.com/frustra/fracture/protobuf"
)

type GameConnection struct {
	Server        *Server
	Conn          net.Conn
	ConnEncrypted *protocol.AESConn

	Player       *protobuf.Player
	EntityServer *network.InternalConnection
}

func (cc *GameConnection) HandleEncryptedConnection() {
	defer func() {
		cc.EntityServer.SendMessage(&protobuf.PlayerAction{
			Player: cc.Player,
			Action: protobuf.PlayerAction_LEAVE,
		})
	}()
	go func() {
		for cc.EntityServer != nil {
			time.Sleep(1 * time.Second)
			cc.Server.PlayerConnections[cc.Player.Uuid] <- protocol.CreatePacket(protocol.KeepAliveID, int32(time.Now().Nanosecond()))
		}
	}()
	for cc.EntityServer != nil {
		cc.Conn.SetReadDeadline(time.Now().Add(time.Second * 30))
		id, buf, err := protocol.ReadPacket(cc.ConnEncrypted)
		n := 0

		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading packet: %s", err.Error())
			}
			return
		} else if id == 0x01 {
			message, _ := protocol.ReadString(buf, 0)
			if len(message) > 0 && message[0] != '/' {
				cc.EntityServer.SendMessage(&protobuf.ChatMessage{
					Message: message,
					Uuid:    cc.Player.Uuid,
				})
			}
		} else if id == 0x04 {
			cc.Player.X, n = protocol.ReadDouble(buf, n)
			cc.Player.FeetY, n = protocol.ReadDouble(buf, n)
			cc.Player.HeadY, n = protocol.ReadDouble(buf, n)
			cc.Player.Z, n = protocol.ReadDouble(buf, n)
			cc.Player.OnGround, n = protocol.ReadBool(buf, n)

			cc.EntityServer.SendMessage(&protobuf.PlayerAction{
				Player: cc.Player,
				Action: protobuf.PlayerAction_MOVE_ABSOLUTE,
				Flags:  1,
			})
		} else if id == 0x05 {
			cc.Player.Yaw, n = protocol.ReadFloat(buf, n)
			cc.Player.Pitch, n = protocol.ReadFloat(buf, n)
			cc.Player.OnGround, n = protocol.ReadBool(buf, n)

			cc.EntityServer.SendMessage(&protobuf.PlayerAction{
				Player: cc.Player,
				Action: protobuf.PlayerAction_MOVE_ABSOLUTE,
				Flags:  2,
			})
		} else if id == 0x06 {
			cc.Player.X, n = protocol.ReadDouble(buf, n)
			cc.Player.FeetY, n = protocol.ReadDouble(buf, n)
			cc.Player.HeadY, n = protocol.ReadDouble(buf, n)
			cc.Player.Z, n = protocol.ReadDouble(buf, n)
			cc.Player.Yaw, n = protocol.ReadFloat(buf, n)
			cc.Player.Pitch, n = protocol.ReadFloat(buf, n)
			cc.Player.OnGround, n = protocol.ReadBool(buf, n)

			cc.EntityServer.SendMessage(&protobuf.PlayerAction{
				Player: cc.Player,
				Action: protobuf.PlayerAction_MOVE_ABSOLUTE,
				Flags:  3,
			})
		} else if id == 0x07 {
			status, n := protocol.ReadByte(buf, n)
			x, n := protocol.ReadInt(buf, n)
			y, n := protocol.ReadByte(buf, n)
			z, n := protocol.ReadInt(buf, n)
			face, n := protocol.ReadByte(buf, n)

			chunkServer, err := cc.Server.FindChunkServer(chunk.WorldCoordsToChunk(int64(x), int64(z)))
			if err != nil {
				log.Print("Tried to destroy block on missing chunk server: ", err)
			}

			chunkServer.SendMessage(&protobuf.BlockUpdate{
				X:       int64(x),
				Y:       uint32(y),
				Z:       int64(z),
				Destroy: true,
				Uuid:    cc.Player.Uuid,
			})
			log.Printf("digging block %d, %d, %d - status %d - face %d", x, y, z, status, face)
		}
	}
}

func (cc *GameConnection) HandleNewConnection() {
	defer func() {
		delete(cc.Server.PlayerConnections, cc.Player.Uuid)
		cc.EntityServer = nil
	}()

	cc.Server.PlayerConnections[cc.Player.Uuid] = make(chan *protocol.Packet, 256)
	cc.Player.EntityId = int64(len(cc.Server.PlayerConnections)+1) + cc.Server.Offset
	cc.Player.HeadY = 105
	cc.Player.FeetY = 105 - 1.62
	cc.Player.OnGround = true

	var err error
	cc.EntityServer, err = cc.Server.FindEntityServer(cc.Player)
	if err != nil {
		log.Printf("Failed to connect to entity server: %s", err)
		protocol.WriteNewPacket(cc.ConnEncrypted, protocol.PreAuthKickID, protocol.CreateJsonMessage("Failed to connect to entity server!", ""))
		return
	}
	cc.EntityServer.SendMessage(&protobuf.PlayerAction{
		Player: cc.Player,
		Action: protobuf.PlayerAction_JOIN,
	})

	var x, z int64
	for x = -8; x < 8; x++ {
		for z = -8; z < 8; z++ {
			conn, err := cc.Server.FindChunkServer(x, z)
			if err != nil {
				log.Printf("Failed to connect to chunk server (%d, %d): %s", x, z, err)
				protocol.WriteNewPacket(cc.ConnEncrypted, protocol.PreAuthKickID, protocol.CreateJsonMessage("Failed to connect to chunk server!", ""))
				return
			}
			conn.SendMessage(&protobuf.ChunkRequest{
				X:    x,
				Z:    z,
				Uuid: cc.Player.Uuid,
			})
		}
	}

	protocol.WriteNewPacket(cc.ConnEncrypted, protocol.LoginSuccessID, cc.Player.Uuid, cc.Player.Username)
	protocol.WriteNewPacket(cc.ConnEncrypted, protocol.JoinGameID, int32(cc.Player.EntityId), uint8(1), byte(0), uint8(1), uint8(cc.Server.Size), "default")
	protocol.WriteNewPacket(cc.ConnEncrypted, protocol.SpawnPositionID, int32(0), int32(128), int32(0))
	protocol.WriteNewPacket(cc.ConnEncrypted, protocol.PlayerAbilitiesID, byte(1), float32(0.05), float32(0.1))
	protocol.WriteNewPacket(cc.ConnEncrypted, protocol.PlayerPositionAndLookID, cc.Player.X, cc.Player.HeadY, cc.Player.Z, cc.Player.Yaw, cc.Player.Pitch, cc.Player.OnGround)

	go cc.Server.DrainPlayerConnections(cc)
	cc.HandleEncryptedConnection()
}
