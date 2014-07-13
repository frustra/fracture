package network

import (
	"encoding/binary"
	"io"
	"log"
	"net"

	"code.google.com/p/gogoprotobuf/proto"
	"github.com/frustra/fracture/protobuf"
)

type MessageHandler interface {
	HandleMessage(message interface{}, conn *InternalConnection)
}

type InternalConnection struct {
	handler MessageHandler
	conn    net.Conn
	queue   chan []byte
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

		c.handler.HandleMessage(message.GetValue(), c)
	}
}

func (c *InternalConnection) FlushQueue() error {
	sizeBuffer := make([]byte, 10)

	for {
		buf := <-c.queue
		if buf == nil {
			break
		}

		l := binary.PutUvarint(sizeBuffer, uint64(len(buf)))
		if _, err := c.conn.Write(sizeBuffer[0:l]); err != nil {
			return err
		}
		if _, err := c.conn.Write(buf); err != nil {
			return err
		}
	}
	return nil
}

func (c *InternalConnection) SendMessage(val interface{}) error {
	message := &protobuf.InternalMessage{}
	message.SetValue(val)

	buf, err := proto.Marshal(message)
	if err != nil {
		return err
	}

	c.queue <- buf
	return nil
}

func (c *InternalConnection) Close() error {
	c.queue <- nil
	return c.conn.Close()
}

func (c *InternalConnection) String() string {
	return c.conn.RemoteAddr().String()
}

func StartInternalConnection(conn net.Conn, handler MessageHandler) *InternalConnection {
	i := &InternalConnection{
		handler: handler,
		conn:    conn,
		queue:   make(chan []byte, 256),
	}

	go func() {
		err := i.FlushQueue()
		if err != nil && err != io.EOF {
			log.Print("Internal outgoing connection error: ", err)
		}
		i.Close()
	}()

	go func() {
		err := i.Handle()
		if err != nil && err != io.EOF {
			log.Print("Internal incoming connection error: ", err)
		}
		i.Close()
	}()

	return i
}

func ServeInternal(addr string, handler MessageHandler) error {
	listener, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Error accepting connection: ", err)
			continue
		}

		StartInternalConnection(conn, handler)
	}
}

func ConnectInternal(addr string, handler MessageHandler) (*InternalConnection, error) {
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		return nil, err
	}
	return StartInternalConnection(conn, handler), nil
}
