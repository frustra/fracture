package main

import (
	"log"

	"github.com/frustra/fracture/network"
	"github.com/frustra/fracture/protobuf"
)

type H struct {
}

func (s *H) HandleMessage(message interface{}, conn *network.InternalConnection) {
	log.Print("Handler invoked: ", message)
}

func main() {
	i, err := network.ConnectInternal("127.0.0.1:25565", &H{})
	if err != nil {
		panic(err)
	}
	var x, z int64 = 1, 3
	i.SendMessage(&protobuf.ChunkRequest{
		X: &x,
		Z: &z,
	})
	var c chan bool
	c <- true
}
