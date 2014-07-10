package edge

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/frustra/fracture/edge/protocol"
	"github.com/frustra/fracture/network"
)

type Server struct {
	Addr       string
	MaxPlayers int
	Cluster    *network.Cluster

	keyPair *protocol.KeyPair

	Clients map[*ClientConnection]bool
}

type ClientConnection struct {
	Server        *Server
	Conn          *net.TCPConn
	ConnEncrypted *protocol.AESConn
	Connected     bool

	Username string
}

func (cc *ClientConnection) HandleEncryptedConnection() {
	cc.Connected = true
	defer func() {
		for client, connected := range cc.Server.Clients {
			if !connected {
				protocol.WriteNewPacket(client.ConnEncrypted, 0x38, cc.Username, false, int16(0))
				protocol.WriteNewPacket(client.ConnEncrypted, 0x02, protocol.CreateJsonMessage(cc.Username+" left the game", "yellow"))
			}
		}
	}()
	go func() {
		for cc.Connected {
			time.Sleep(1 * time.Second)
			protocol.WriteNewPacket(cc.ConnEncrypted, 0x00, int32(time.Now().Nanosecond()))
		}
	}()
	for cc.Connected {
		cc.Conn.SetReadDeadline(time.Now().Add(time.Second * 30))
		id, buf, err := protocol.ReadPacket(cc.ConnEncrypted)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error: %s\n", err.Error())
			}
			return
		} else if id == 0x01 {
			message, _ := protocol.ReadString(buf, 0)
			if len(message) > 0 && message[0] != '/' {
				for client, _ := range cc.Server.Clients {
					if client.Connected {
						protocol.WriteNewPacket(client.ConnEncrypted, 0x02, protocol.CreateJsonMessage("<"+cc.Username+"> "+message, ""))
					}
				}
			}
		}
	}
}

func (cc *ClientConnection) HandleConnection() {
	defer func() {
		delete(cc.Server.Clients, cc)
		cc.Connected = false
		cc.Conn.Close()
	}()
	remoteAddr := cc.Conn.RemoteAddr().String()

	state := 0
	verifyToken := make([]byte, 0)
	for cc.Connected {
		cc.Conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		id, buf, err := protocol.ReadPacket(cc.Conn)
		if err != nil {
			err2, ok := err.(net.Error)
			if ok && err2.Timeout() {
				fmt.Printf("Error Timeout: %s\n", err2.Error())
				return
			} else {
				return
			}
		} else {
			switch id {
			case 0x00:
				if state == 1 {
					fmt.Printf("Server pinged from: %s\n", remoteAddr)
					protocol.WriteNewPacket(cc.Conn, 0x00, protocol.CreateStatusResponse("1.7.10", 5, 0, cc.Server.MaxPlayers, protocol.CreateJsonMessage("Fracture Distributed Server", "green")))
				} else if state == 2 {
					cc.Username, _ = protocol.ReadString(buf, 0)
					fmt.Printf("Got connection from %s\n", cc.Username)
					defer fmt.Printf("Connection closed for %s\n", cc.Username)

					pubKey := cc.Server.keyPair.Serialize()
					verifyToken = protocol.GenerateKey(16)
					protocol.WriteNewPacket(cc.Conn, 0x01, "", int16(len(pubKey)), pubKey, int16(len(verifyToken)), verifyToken)
				} else {
					_, n := protocol.ReadUvarint(buf, 0) // version
					_, n = protocol.ReadString(buf, n)   // address
					_, n = protocol.ReadShort(buf, n)    // port
					nextstate, n := protocol.ReadUvarint(buf, n)
					state = int(nextstate)
				}
			case 0x01:
				if state == 2 {
					secretLen, n := protocol.ReadShort(buf, 0)
					secretEncrypted, n := protocol.ReadBytes(buf, n, secretLen)
					tokenLen, n := protocol.ReadShort(buf, n)
					tokenEncrypted, n := protocol.ReadBytes(buf, n, tokenLen)

					verifyToken2, err := protocol.DecryptRSABytes(tokenEncrypted, cc.Server.keyPair)
					if err != nil {
						fmt.Printf("Error: %s\n", err.Error())
						return
					} else if !bytes.Equal(verifyToken, verifyToken2) {
						fmt.Printf("Error: Verify token did not match!")
						return
					}
					sharedSecret, err := protocol.DecryptRSABytes(secretEncrypted, cc.Server.keyPair)
					if err != nil {
						fmt.Printf("Error: %s\n", err.Error())
						return
					}

					cc.ConnEncrypted, err = protocol.NewAESConn(cc.Conn, sharedSecret)
					if err != nil {
						fmt.Printf("Error: %s\n", err.Error())
					}

					uuid, err := protocol.CheckAuth(cc.Username, "", cc.Server.keyPair, sharedSecret)
					if err != nil {
						fmt.Printf("Failed to verify username: %s\n%s\n", cc.Username, err)
						protocol.WriteNewPacket(cc.ConnEncrypted, 0x00, protocol.CreateJsonMessage("Failed to verify username!", ""))
						return
					}
					protocol.WriteNewPacket(cc.ConnEncrypted, 0x02, uuid, cc.Username)
					protocol.WriteNewPacket(cc.ConnEncrypted, 0x01, int32(1), uint8(0), byte(0), uint8(1), uint8(cc.Server.MaxPlayers), "default")
					protocol.WriteNewPacket(cc.ConnEncrypted, 0x05, int32(0), int32(0), int32(0))
					protocol.WriteNewPacket(cc.ConnEncrypted, 0x39, byte(1), float32(5), float32(5))
					protocol.WriteNewPacket(cc.ConnEncrypted, 0x08, float64(0), float64(128), float64(0), float32(0), float32(0), false)
					protocol.WriteNewPacket(cc.ConnEncrypted, 0x26, int16(0), int32(0), true)
					cc.HandleEncryptedConnection()
					// protocol.WriteNewPacket(cc.ConnEncrypted, 0x00, protocol.CreateJsonMessage("Server will bbl", "blue"))
					return
				} else {
					time, _ := protocol.ReadLong(buf, 0)
					//fmt.Printf("Ping: %d\n", time)
					protocol.WriteNewPacket(cc.Conn, 0x01, time)
				}
			default:
				fmt.Printf("Unknown Packet (state:%d): 0x%X : %s\n", state, id, hex.Dump(buf))
			}
		}
	}
}

func (s *Server) Serve() {
	addr, err := net.ResolveTCPAddr("tcp4", s.Addr)
	assertNoErr(err)

	listener, err := net.ListenTCP("tcp", addr)
	assertNoErr(err)

	tmpkey, err := protocol.GenerateKeyPair(1024)
	assertNoErr(err)
	s.keyPair = tmpkey

	fmt.Printf("Listening for TCP on %s\n", s.Addr)
	defer listener.Close()

	s.Clients = make(map[*ClientConnection]bool)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}

		client := &ClientConnection{s, conn, nil, true, ""}
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

func assertNoErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err.Error())
		os.Exit(1)
	}
}
