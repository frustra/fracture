package network

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/gogoprotobuf/proto"
	"github.com/frustra/fracture/protobuf"
	"github.com/frustra/fracture/world"
	"github.com/hashicorp/memberlist"
)

const (
	EdgeType uint32 = iota
	EntityType
	ChunkType
)

var (
	ErrMissingChunk = errors.New("missing chunk server")
)

type Node struct {
	Meta       *protobuf.NodeMeta
	Connection *InternalConnection
}

type NodeMap map[string]*Node
type ChunkNodeMap1D map[int64]*Node
type ChunkNodeMap map[int64]ChunkNodeMap1D

type Cluster struct {
	Members *memberlist.Memberlist
	config  *memberlist.Config

	existing string

	LocalNodeMeta  *protobuf.NodeMeta
	nodeMetaBuffer []byte

	NodesByName NodeMap
	EdgeNodes   NodeMap
	EntityNodes NodeMap
	ChunkNodes  ChunkNodeMap
}

func CreateCluster(addr, existing string) (*Cluster, error) {
	host, portString, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(portString)
	if err != nil || port == 0 {
		return nil, errors.New("port must be defined")
	}

	config := memberlist.DefaultLocalConfig()
	config.BindPort = port
	config.AdvertisePort = port

	if host != "" {
		config.BindAddr = host
	}

	c := &Cluster{
		config:   config,
		existing: existing,

		LocalNodeMeta: &protobuf.NodeMeta{},

		NodesByName: make(NodeMap),
		EdgeNodes:   make(NodeMap),
		EntityNodes: make(NodeMap),
		ChunkNodes:  make(ChunkNodeMap),
	}

	config.Delegate = c
	config.Events = c
	config.LogOutput = ioutil.Discard
	config.Name = uuid.New()

	return c, nil
}

func (c *Cluster) Name() string {
	return c.Members.LocalNode().Name
}

func (c *Cluster) Join() error {
	nodeMeta, err := proto.Marshal(c.LocalNodeMeta)
	if err != nil {
		return err
	}

	c.nodeMetaBuffer = nodeMeta

	m, err := memberlist.Create(c.config)
	if err != nil {
		return err
	}

	c.Members = m
	_, err = m.Join([]string{c.existing})
	return err
}

func (c *Cluster) Part() {
	c.Members.Leave(2 * time.Second)
	c.Members.Shutdown()
}

func (c *Cluster) ChunkConnection(x, z int64, h MessageHandler) (*InternalConnection, error) {
	x, z = world.ChunkCoordsToNode(x, z)

	zm := c.ChunkNodes[z]
	if zm == nil {
		return nil, ErrMissingChunk
	}

	node := zm[x]
	if node == nil {
		return nil, ErrMissingChunk
	}

	if node.Connection != nil {
		return node.Connection, nil
	}

	conn, err := ConnectInternal(node.Meta.Addr, h)
	if err != nil {
		return nil, err
	}

	node.Connection = conn
	return conn, nil
}

func (c *Cluster) NotifyJoin(n *memberlist.Node) {
	meta := &protobuf.NodeMeta{}
	err := proto.Unmarshal(n.Meta, meta)
	if err != nil {
		log.Print("error unmarshalling node metadata: ", err)
		return
	}

	metaHost, metaPort, err := net.SplitHostPort(meta.Addr)
	if metaHost == "" {
		meta.Addr = n.Addr.String() + ":" + metaPort
	}

	node := &Node{Meta: meta}

	if _, exists := c.NodesByName[n.Name]; exists {
		log.Printf("Duplicate node %s joined from %s:%d", n.Name, n.Addr, n.Port)
		return
	}

	c.NodesByName[n.Name] = node

	switch meta.Type {
	case EdgeType:
		c.EdgeNodes[n.Name] = node
	case EntityType:
		c.EntityNodes[n.Name] = node
	case ChunkType:
		x, z := meta.GetX(), meta.GetZ()
		if c.ChunkNodes[z] == nil {
			c.ChunkNodes[z] = make(ChunkNodeMap1D)
		}
		c.ChunkNodes[z][x] = node
	default:
		log.Printf("Unknown node type %d joined from %s:%d", meta.Type, n.Addr, n.Port)
		return
	}

	// log.Printf("%s %s:%s [%s] joined from %s:%d", meta.GetType(), metaHost, metaPort, n.Name, n.Addr, n.Port)
}

func (c *Cluster) NotifyLeave(n *memberlist.Node) {
	node, exists := c.NodesByName[n.Name]

	if exists {
		switch node.Meta.Type {
		case EdgeType:
			delete(c.EdgeNodes, n.Name)
		case EntityType:
			delete(c.EntityNodes, n.Name)
		case ChunkType:
			x, z := node.Meta.GetX(), node.Meta.GetZ()

			zm := c.ChunkNodes[z]
			if zm != nil {
				delete(zm, x)
				if len(zm) == 0 {
					delete(c.ChunkNodes, z)
				}
			} else {
				log.Printf("Unmapped chunk node %s left from %s:%d", n.Name, n.Addr, n.Port)
			}
		}

		delete(c.NodesByName, n.Name)
	} else {
		log.Printf("Unknown node %s left from %s:%d", n.Name, n.Addr, n.Port)
	}
}

func (c *Cluster) NotifyUpdate(n *memberlist.Node) {
	node := c.NodesByName[n.Name]
	if node == nil {
		log.Printf("Unknown node %s updated from %s:%d", n.Name, n.Addr, n.Port)
	}

	meta := &protobuf.NodeMeta{}
	err := proto.Unmarshal(n.Meta, meta)
	if err != nil {
		log.Print("error unmarshalling node metadata: ", err)
		return
	}

	node.Meta = meta
}

func (c *Cluster) NodeMeta(limit int) []byte {
	return c.nodeMetaBuffer
}

func (c *Cluster) NotifyMsg([]byte) {

}

func (c *Cluster) GetBroadcasts(overhead, limit int) [][]byte {
	return nil
}

func (c *Cluster) LocalState(join bool) []byte {
	return nil
}

func (c *Cluster) MergeRemoteState(buf []byte, join bool) {

}
