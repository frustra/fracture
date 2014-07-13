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
	"github.com/hashicorp/memberlist"
)

type Node struct {
	Meta       *protobuf.NodeMeta
	Connection *InternalConnection
}

type NodeMap map[string]*Node

type Cluster struct {
	Members *memberlist.Memberlist
	config  *memberlist.Config

	existing string

	LocalNodeMeta  protobuf.NodeMeta
	nodeMetaBuffer []byte

	Edge, Entity, Chunk NodeMap
	TypeLookup          map[string]NodeMap
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
		config:     config,
		existing:   existing,
		Edge:       make(NodeMap),
		Entity:     make(NodeMap),
		Chunk:      make(NodeMap),
		TypeLookup: make(map[string]NodeMap),
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
	nodeMeta, err := proto.Marshal(&c.LocalNodeMeta)
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

	var m NodeMap

	switch meta.Type {
	case "edge":
		m = c.Edge
	case "entity":
		m = c.Entity
	case "chunk":
		m = c.Chunk
	default:
		log.Printf("Unknown node type %s joined from %s:%d", meta.Type, n.Addr, n.Port)
		return
	}

	c.TypeLookup[n.Name] = m

	if _, exists := m[n.Name]; exists {
		log.Printf("Duplicate %s node %s joined from %s:%d", meta.Type, n.Name, n.Addr, n.Port)
		return
	}

	m[n.Name] = &Node{Meta: meta}

	// log.Printf("%s %s:%s [%s] joined from %s:%d", meta.GetType(), metaHost, metaPort, n.Name, n.Addr, n.Port)
}

func (c *Cluster) NotifyLeave(n *memberlist.Node) {
	typeMap, ok := c.TypeLookup[n.Name]
	if ok {
		delete(typeMap, n.Name)
	}
	delete(c.TypeLookup, n.Name)

	// log.Printf("%s left", n.Name)
}

func (c *Cluster) NotifyUpdate(n *memberlist.Node) {
	typeMap := c.TypeLookup[n.Name]

	meta := &protobuf.NodeMeta{}
	err := proto.Unmarshal(n.Meta, meta)
	if err != nil {
		log.Print("error unmarshalling node metadata: ", err)
		return
	}
	typeMap[n.Name].Meta = meta
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
