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
	}

	config.Delegate = c
	config.Events = c
	config.LogOutput = ioutil.Discard
	config.Name = uuid.New()

	return c, nil
}

type Cluster struct {
	Members *memberlist.Memberlist
	config  *memberlist.Config

	existing string

	NodeType, NodeAddr string
	nodeMeta           []byte
}

func (c *Cluster) Name() string {
	return c.Members.LocalNode().Name
}

func (c *Cluster) Join() error {
	meta := &protobuf.NodeMeta{
		Addr: proto.String(c.NodeAddr),
		Type: proto.String(c.NodeType),
	}

	nodeMeta, err := proto.Marshal(meta)
	if err != nil {
		return err
	}

	c.nodeMeta = nodeMeta

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
	}

	metaHost, metaPortString, err := net.SplitHostPort(meta.GetAddr())
	if metaHost == "" {
		metaHost = n.Addr.String()
	}

	log.Printf("%s %s:%s [%s] joined from %s:%d", meta.GetType(), metaHost, metaPortString, n.Name, n.Addr, n.Port)
}

func (c *Cluster) NotifyLeave(n *memberlist.Node) {
	log.Printf("%s left", n.Name)
}

func (c *Cluster) NotifyUpdate(n *memberlist.Node) {

}

func (c *Cluster) NodeMeta(limit int) []byte {
	return c.nodeMeta
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
