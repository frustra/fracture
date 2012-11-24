package edge

import (
	"net"
	"net/http"
	"net/rpc"
	"log"
	"fracture/master"
)

type ChunkNode struct {
	X, Y int64
	Addr string
	client *rpc.Client
}

type ChunkNodes struct {
	nodes map[string]*ChunkNode
}

func (n *ChunkNodes) AddChunkNode(c *ChunkNode, unused *bool) error {
	n.nodes[c.Addr] = c
	var err error
	c.client, err = rpc.DialHTTP("tcp", c.Addr)
	if err != nil {
		log.Print(err)
	}
	log.Printf("Chunk server %s connected", c.Addr)
	return nil
}

func (n *ChunkNodes) RemoveChunkNode(c *ChunkNode, unused *bool) error {
	cc, exists := n.nodes[c.Addr]
	if !exists {
		return nil
	}
	delete(n.nodes, cc.Addr)
	cc.client.Close()
	log.Printf("Chunk server %s disconnected", cc.Addr)
	return nil
}

var chunkNodes *ChunkNodes

func ListenMaster() {
	chunkNodes = new(ChunkNodes)
	chunkNodes.nodes = make(map[string]*ChunkNode)

	rpc.Register(chunkNodes)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting up master connection on %s", l.Addr().String())

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:25566")
	if err != nil {
		log.Fatal(err)
	}
	client.Call("Server.AnnounceEdgeServer", &master.EdgeServer{Addr: l.Addr().String()}, nil)
	client.Close()

	http.Serve(l, nil)
}
