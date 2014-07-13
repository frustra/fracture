package edge

import (
	"bytes"
	"encoding/hex"
	"log"
	"net"
	"time"

	"github.com/frustra/fracture/edge/protocol"
)

func (cc *GameConnection) HandleAuth() {
	remoteAddr := cc.Conn.RemoteAddr().String()

	state := 0
	verifyToken := make([]byte, 0)
	for {
		cc.Conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		id, buf, err := protocol.ReadPacket(cc.Conn)
		if err != nil {
			err2, ok := err.(net.Error)
			if ok && err2.Timeout() {
				log.Printf("Timeout handling connection (%s): %s", remoteAddr, err2.Error())
			}
			return
		} else {
			switch id {
			case 0x00: // Handshake, Status Request, Login start
				if state == 1 {
					log.Printf("Server pinged from: %s", remoteAddr)

					protocol.WriteNewPacket(cc.Conn, protocol.StatusResponseID, cc.Server.GetMinecraftStatus())
				} else if state == 2 {
					cc.Player.Username, _ = protocol.ReadString(buf, 0)
					log.Printf("Got connection from %s", cc.Player.Username)
					defer log.Printf("Connection closed for %s", cc.Player.Username)

					pubKey := cc.Server.keyPair.Serialize()
					verifyToken = protocol.GenerateKey(16)
					protocol.WriteNewPacket(cc.Conn, protocol.EncryptionRequestID, "", int16(len(pubKey)), pubKey, int16(len(verifyToken)), verifyToken)
				} else {
					_, n := protocol.ReadVarint(buf, 0) // version
					_, n = protocol.ReadString(buf, n)  // address
					_, n = protocol.ReadShort(buf, n)   // port

					nextstate, n := protocol.ReadVarint(buf, n)
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

					// protocol.WriteNewPacket(cc.ConnEncrypted, protocol.PreAuthKickID, protocol.CreateJsonMessage("Server will bbl", "blue"))
					cc.HandleNewConnection()
					return
				} else {
					time, _ := protocol.ReadLong(buf, 0)
					//fmt.Printf("Ping: %d\n", time)
					protocol.WriteNewPacket(cc.Conn, protocol.PingResponseID, time)
					return
				}
			default:
				log.Printf("Unknown Packet (state:%d): 0x%X : %s", state, id, hex.Dump(buf))
			}
		}
	}
}
