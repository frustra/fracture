package edge

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"io"
	"log"
	"net"
	"time"

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
	var err error
	cc.EntityServer, err = cc.Server.FindEntityServer(cc.Player)
	if err != nil {
		log.Printf("Failed to connect to entity server: %s", err)
		protocol.WriteNewPacket(cc.ConnEncrypted, protocol.DisconnectID, protocol.CreateJsonMessage("Failed to connect to entity server!", ""))
		return
	}
	cc.EntityServer.SendMessage(&protobuf.PlayerAction{
		Player: cc.Player,
		Action: protobuf.PlayerAction_JOIN,
	})
	defer func() {
		for client, _ := range cc.Server.Clients {
			if client.EntityServer != nil {
				protocol.WriteNewPacket(client.ConnEncrypted, protocol.PlayerListItemID, cc.Player.Username, false, int16(0))
				protocol.WriteNewPacket(client.ConnEncrypted, protocol.ChatMessageID, protocol.CreateJsonMessage(cc.Player.Username+" left the game", "yellow"))
			}
		}
	}()
	go func() {
		for cc.EntityServer != nil {
			time.Sleep(1 * time.Second)
			protocol.WriteNewPacket(cc.ConnEncrypted, protocol.KeepAliveID, int32(time.Now().Nanosecond()))
		}
	}()
	for cc.EntityServer != nil {
		cc.Conn.SetReadDeadline(time.Now().Add(time.Second * 30))
		id, buf, err := protocol.ReadPacket(cc.ConnEncrypted)
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading packet: %s", err.Error())
			}
			return
		} else if id == 0x01 {
			message, _ := protocol.ReadString(buf, 0)
			if len(message) > 0 && message[0] != '/' {
				for client, _ := range cc.Server.Clients {
					if client.EntityServer != nil {
						protocol.WriteNewPacket(client.ConnEncrypted, protocol.ChatMessageID, protocol.CreateJsonMessage("<"+cc.Player.Username+"> "+message, ""))
					}
				}
			}
		}
	}
}

func (cc *GameConnection) HandleConnection() {
	defer func() {
		delete(cc.Server.Clients, cc)
		cc.EntityServer = nil
		cc.Conn.Close()
	}()
	remoteAddr := cc.Conn.RemoteAddr().String()

	state := 0
	verifyToken := make([]byte, 0)
	for {
		cc.Conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		id, buf, err := protocol.ReadPacket(cc.Conn)
		if err != nil {
			err2, ok := err.(net.Error)
			if ok && err2.Timeout() {
				log.Printf("Timeout handling connection: %s", err2.Error())
				return
			} else {
				return
			}
		} else {
			switch id {
			case 0x00: // Handshake, Status Request, Login Start
				if state == 1 {
					log.Printf("Server pinged from: %s", remoteAddr)
					protocol.WriteNewPacket(cc.Conn, protocol.StatusResponseID, protocol.CreateStatusResponse("1.7.10", 5, 0, cc.Server.Size, protocol.CreateJsonMessage("Fracture Distributed Server", "green")))
				} else if state == 2 {
					cc.Player.Username, _ = protocol.ReadString(buf, 0)
					log.Printf("Got connection from %s", cc.Player.Username)
					defer log.Printf("Connection closed for %s", cc.Player.Username)

					pubKey := cc.Server.keyPair.Serialize()
					verifyToken = protocol.GenerateKey(16)
					protocol.WriteNewPacket(cc.Conn, protocol.EncryptionRequestID, "", int16(len(pubKey)), pubKey, int16(len(verifyToken)), verifyToken)
				} else {
					_, n := protocol.ReadUvarint(buf, 0) // version
					_, n = protocol.ReadString(buf, n)   // address
					_, n = protocol.ReadShort(buf, n)    // port
					nextstate, n := protocol.ReadUvarint(buf, n)
					state = int(nextstate)
				}
			case 0x01: // Encryption Response, Ping Request
				if state == 2 {
					secretLen, n := protocol.ReadShort(buf, 0)
					secretEncrypted, n := protocol.ReadBytes(buf, n, secretLen)
					tokenLen, n := protocol.ReadShort(buf, n)
					tokenEncrypted, n := protocol.ReadBytes(buf, n, tokenLen)

					verifyToken2, err := protocol.DecryptRSABytes(tokenEncrypted, cc.Server.keyPair)
					if err != nil {
						log.Printf("Error decrypting RSA token: %s", err.Error())
						return
					} else if !bytes.Equal(verifyToken, verifyToken2) {
						log.Printf("Error: verification token did not match!")
						return
					}
					sharedSecret, err := protocol.DecryptRSABytes(secretEncrypted, cc.Server.keyPair)
					if err != nil {
						log.Printf("Error decrypting RSA secret: %s", err.Error())
						return
					}

					cc.ConnEncrypted, err = protocol.NewAESConn(cc.Conn, sharedSecret)
					if err != nil {
						log.Printf("Error creating AES connection: %s", err.Error())
					}

					cc.Player.Uuid, err = protocol.CheckAuth(cc.Player.Username, "", cc.Server.keyPair, sharedSecret)
					if err != nil {
						log.Printf("Failed to verify username %s: %s", cc.Player.Username, err)
						protocol.WriteNewPacket(cc.ConnEncrypted, protocol.PreAuthKickID, protocol.CreateJsonMessage("Failed to verify username!", ""))
						return
					}
					protocol.WriteNewPacket(cc.ConnEncrypted, protocol.LoginSuccessID, cc.Player.Uuid, cc.Player.Username)
					protocol.WriteNewPacket(cc.ConnEncrypted, protocol.JoinGameID, int32(1), uint8(0), byte(0), uint8(1), uint8(cc.Server.Size), "default")
					protocol.WriteNewPacket(cc.ConnEncrypted, protocol.SpawnPositionID, int32(0), int32(0), int32(0))
					protocol.WriteNewPacket(cc.ConnEncrypted, protocol.PlayerAbilitiesID, byte(1), float32(5), float32(5))
					protocol.WriteNewPacket(cc.ConnEncrypted, protocol.PlayerPositionAndLookID, float64(0), float64(128), float64(0), float32(0), float32(0), false)
					worldData := make([]byte, 4096+2048+2048+2048+256)
					for i := 0; i < 4096; i++ {
						if i >= 3840 {
							worldData[i] = 2
						} else {
							worldData[i] = 3
						}
					}
					var compressed bytes.Buffer

					w := zlib.NewWriter(&compressed)
					for x := -10; x <= 10; x++ {
						for y := -10; y <= 10; y++ {
							w.Write(worldData)
						}
					}
					w.Close()
					chunkData := []interface{}{int16(21 * 21), int32(compressed.Len()), true, compressed.Bytes()}
					for x := -10; x <= 10; x++ {
						for y := -10; y <= 10; y++ {
							chunkData = append(chunkData, int32(x), int32(y), uint16(1), uint16(0))
						}
					}
					protocol.WriteNewPacket(cc.ConnEncrypted, protocol.MapChunkBulkID, chunkData...)
					cc.HandleEncryptedConnection()
					// protocol.WriteNewPacket(cc.ConnEncrypted, protocol.PreAuthKickID, protocol.CreateJsonMessage("Server will bbl", "blue"))
					return
				} else {
					time, _ := protocol.ReadLong(buf, 0)
					//fmt.Printf("Ping: %d\n", time)
					protocol.WriteNewPacket(cc.Conn, protocol.PingResponseID, time)
				}
			default:
				log.Printf("Unknown Packet (state:%d): 0x%X : %s", state, id, hex.Dump(buf))
			}
		}
	}
}
