package edge

import (
	"flag"
	"log"
	//"fracture/chunk"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"net"
	"strconv"
)

var port int
var port2 int
var serverInfo []byte
var players []string
var privkey *rsa.PrivateKey
var pubkey []byte
var sharedkey []byte
var clientcfb cipher.Stream
var servercfb cipher.Stream

const ENCRYPTION_ENABLED = false

func doHandshake(conn net.Conn) {
	log.Printf("Got connection from %s", conn.RemoteAddr().String())

	buf := make([]byte, 65536)
	encrypted := false
	for {
		length, err := conn.Read(buf)
		if err != nil {
			break
		}
		if encrypted {
			log.Printf("Client: %s", hex.Dump(buf[0:length]))
			decrypt(clientcfb, buf[0:length])
		}
		name := packName[buf[0]]
		//log.Printf("Client: %s", hex.Dump(buf[0:length]))
		switch name {
		case "listping":
			if buf[1] == 0x01 {
				//log.Printf("Server Info: %s", hex.Dump(serverInfo))
				conn.Write(serverInfo)
			}
		case "handshake":
			if buf[1] == 49 { // Protocol version
				n := 2
				username, n := ReadString(buf, n)
				host, n := ReadString(buf, n)
				port, n := ReadInt(buf, n)
				log.Printf("%s, %s, %d", username, host, port)
				if ENCRYPTION_ENABLED {
					packet := CreatePacket("encryptrequest", "555555", int16(len(pubkey)), pubkey, int16(4), []byte{0x05, 0x52, 0x88, 0x04})
					conn.Write(packet.Serialize())
				} else {
					packet := CreatePacket("login", int32(1298), "flat", byte(1), byte(0), byte(0), byte(0), byte(50))
					conn.Write(packet.Serialize())
					packet = CreatePacket("spawnpos", int32(0), int32(32), int32(0))
					conn.Write(packet.Serialize())
					packet = CreatePacket("playerposlook", float64(0), float64(64), float64(0), float64(0), float32(0), float32(0), bool(true))
					conn.Write(packet.Serialize())
				}
			}
		case "encryptresponse":
			if ENCRYPTION_ENABLED {
				n := 1
				keylen, n := ReadShort(buf, n)
				key, n := ReadBytes(buf, n, keylen)
				sharedkey, _ = rsa.DecryptPKCS1v15(rand.Reader, privkey, key)
				clientcipher, _ := aes.NewCipher(sharedkey)
				clientcfb = cipher.NewCFBDecrypter(clientcipher, sharedkey)
				servercipher, _ := aes.NewCipher(sharedkey)
				servercfb = cipher.NewCFBEncrypter(servercipher, sharedkey)
				encrypted = true
				packet := CreatePacket("encryptresponse", int16(0), []byte{}, int16(0), []byte{})
				conn.Write(packet.Serialize())
			}
		case "clientstatuses":
			if ENCRYPTION_ENABLED {
				packet := CreatePacket("login", int32(1298), "flat", byte(1), byte(0), byte(0), byte(0), byte(50))
				tmp := packet.Serialize()
				tmp2 := encrypt(servercfb, tmp)
				log.Printf("Key: %s", hex.Dump(sharedkey))
				log.Printf("Client: %s", hex.Dump(buf[0:length]))
				log.Printf("Server: %s", hex.Dump(tmp))
				log.Printf("Serverc: %s", hex.Dump(tmp2))
				conn.Write(tmp2)
			}
		case "player":
		case "playerpos":
		case "playerlook":
		case "playerposlook":
		case "keepalive":
		case "kick":
			msg, _ := ReadString(buf, 1)
			log.Printf("Client disconnected with message: %s", msg)
		default:
			log.Printf("Unknown: %s", hex.Dump(buf[0:length]))
		}
		//packet := CreatePacket("kick", "Server will bbl.")
		//conn.Write(packet.Serialize())
	}
	conn.Close()
}

func init() {
	InitPackets()
	privkey, _ = rsa.GenerateKey(rand.Reader, 1024)
	pubkey, _ = x509.MarshalPKIXPublicKey(&privkey.PublicKey)
	flag.IntVar(&port, "port", 25565, "TCP port to listen on")
}

func Serve() {
	c := new(chunk.Client)
	c.Connect("localhost:12444")
	c.SetBlock(0, 0, 0, 2)
	ch := c.GetChunk(0)
	log.Printf("%d should equal 2", ch.Data[0][0][0])

	flag.Parse()
	log.Println("Starting up edge server")
	ListenMaster()
	service := ":" + strconv.Itoa(port)
	addr, err := net.ResolveTCPAddr("tcp4", service)
	assertNoErr(err)

	listener, err := net.ListenTCP("tcp", addr)
	assertNoErr(err)

	log.Printf("Listening for TCP on %s", service)

	packet := CreatePacket("kick", JoinStrings(
		[]byte{0xa7, 0x31}, // Magic characters
		"49",               // Protocol Version (1.4.5)
		"1.4.5",            // Minecraft version
		"Fracture Server",  // MOTD
		"0",                // Online players
		"50",               // Max players
	))
	serverInfo = packet.Serialize()

	players = make([]string, 0)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go doHandshake(conn)
	}
}

func assertNoErr(err error) {
	if err != nil {
		log.Fatalf("Fatal: %s", err.Error())
	}
}

func encrypt(cfb cipher.Stream, plaintext []byte) []byte {
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	return ciphertext
}

func decrypt(cfb cipher.Stream, ciphertext []byte) {
	cfb.XORKeyStream(ciphertext, ciphertext)
}
