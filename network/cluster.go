package network

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"time"

	"code.google.com/p/go-uuid/uuid"
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

	config.Events = c
	config.LogOutput = ioutil.Discard
	config.Name = uuid.New()

	return c, nil
}

type Cluster struct {
	Members *memberlist.Memberlist
	config  *memberlist.Config

	existing string
}

func (c *Cluster) Name() string {
	return c.Members.LocalNode().Name
}

func (c *Cluster) Join() error {
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
	log.Printf("%s joined from %s:%d", n.Name, n.Addr, n.Port)
}

func (c *Cluster) NotifyLeave(n *memberlist.Node) {
	log.Printf("%s left", n.Name)
}

func (c *Cluster) NotifyUpdate(n *memberlist.Node) {

}
