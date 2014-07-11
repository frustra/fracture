package network

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	"code.google.com/p/gogoprotobuf/proto"
	"github.com/frustra/fracture/protobuf"
)

type Server interface {
	Serve() error

	NodeType() string
	NodePort() int
}

type InternalConnection struct {
	server Server
	conn   *net.TCPConn
}

type proxyByteReader struct {
	reader io.Reader
}

func (r *proxyByteReader) ReadByte() (c byte, err error) {
	buf := make([]byte, 1)
	if _, err := r.reader.Read(buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (c *InternalConnection) Handle() error {
	defer c.conn.Close()
	// remoteAddr := c.conn.RemoteAddr().String()
	buf := make([]byte, 65536)

	for {
		// c.conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		length, err := binary.ReadUvarint(&proxyByteReader{c.conn})
		if err != nil {
			return err
		}

		packet := buf[0:length]
		var read uint64

		for read < length {
			n, err := c.conn.Read(packet[read:])
			if err != nil {
				return err
			}
			read += uint64(n)
		}

		message := &protobuf.InternalMessage{}
		err = proto.Unmarshal(packet, message)
		if err != nil {
			return err
		}

		log.Print(message)
	}
}

func ServeInternal(addr string, s Server) error {
	laddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Print("Error accepting connection: ", err)
			continue
		}

		client := &InternalConnection{s, conn}
		go func() {
			err := client.Handle()
			if err != nil {
				log.Print("Internal client error: ", err)
			}
		}()
	}
}
