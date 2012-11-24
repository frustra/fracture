package edge

import (
	"bytes"
	//"encoding/hex"
	"flag"
	"fmt"
	"log"
	//"fracture/chunk"
	"net"
	"os"
	"strconv"
)

var port int
var port2 int
var serverInfo []byte
var responseList map[string]string
var players []string

func handleConn(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("Got connection from %s\n", remoteAddr)

	buf := make([]byte, 65536)
	var connout net.Conn
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		//fmt.Printf("Client: %s", hex.Dump(buf[0:n]))
		if buf[0] == 0xFE && buf[1] == 0x01 {
			//fmt.Printf("Server Info: %s", hex.Dump(serverInfo))
			_, err := conn.Write(serverInfo)
			if err != nil {
				break
			}
		} else if connout == nil {
			packet := CreatePacket("kick", "Server will bbl.")
			_, err := conn.Write(packet.Serialize())
			if err != nil {
				break
			}
		} else {
			_, err = connout.Write(buf[0:n])
			if err != nil {
				break
			}
		}
	}
	conn.Close()
	if connout != nil {
		connout.Close()
	}
}

func init() {
	flag.IntVar(&port, "port", 25565, "TCP port to listen on")
}

func Serve() {
	/*c := new(chunk.Client)
	c.Connect("localhost:12444")
	c.SetBlock(0, 0, 0, 2)
	ch := c.GetChunk(0)
	log.Printf("%d should equal 2", ch.Data[0][0][0])*/

	flag.Parse()
	log.Println("Starting up edge server")
	service := ":" + strconv.Itoa(port)
	addr, err := net.ResolveTCPAddr("tcp4", service)
	assertNoErr(err)

	listener, err := net.ListenTCP("tcp", addr)
	assertNoErr(err)

	fmt.Printf("Listening for TCP on %s\n", service)

	responseList = make(map[string]string)
	responseList["hostname"] = "Fracture Server"
	responseList["gametype"] = "SMP"
	responseList["game_id"] = "MINECRAFT"
	responseList["version"] = "1.4.4"
	responseList["plugins"] = "Plugins"
	responseList["map"] = "lobby"
	responseList["numplayers"] = "0"
	responseList["maxplayers"] = "50"
	responseList["hostport"] = "25565"
	responseList["hostip"] = "63.141.238.132"
	serverInfo = []byte{0xff, 0, 0, 0, 0xa7, 0, 0x31, 0, 0, 0,
		0x34, 0, 0x39, 0, 0} // Protocol version: '49'
	serverInfo = bytes.Join([][]byte{
		serverInfo,
		toUtf16ByteArrayA(responseList["version"]), []byte{0, 0},
		toUtf16ByteArrayA(responseList["hostname"]), []byte{0, 0},
		toUtf16ByteArrayA(responseList["numplayers"]), []byte{0, 0},
		toUtf16ByteArrayA(responseList["maxplayers"])},
		[]byte{})

	serverInfo[3] = (byte)(((len(serverInfo) - 3) >> 9) & 0xFF)
	serverInfo[2] = (byte)(((len(serverInfo) - 3) >> 1) & 0xFF)
	players = make([]string, 0)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleConn(conn)
	}
}

func toUtf16ByteArrayA(str string) []byte {
	buf := make([]byte, len(str)*2)
	for i := 0; i < len(str); i++ {
		buf[i*2] = 0
		buf[i*2+1] = str[i]
	}
	return buf
}

func assertNoErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: %s\n", err.Error())
		os.Exit(1)
	}
}
