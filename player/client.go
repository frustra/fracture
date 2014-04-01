package player

import (
	"log"
	"net/rpc"
)

type Client struct {
	client *rpc.Client
}

func (c *Client) Connect(addr string) {
	var err error
	c.client, err = rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) GetPlayer(id uint) (ret *Player) {
	err := c.client.Call("Group.GetPlayer", id, &ret)
	if err != nil {
		log.Fatalf("error getting player %d: %s", id, err)
	}
	return
}

func (c *Client) GetPlayers() (ret *[16]*Player) {
	err := c.client.Call("Group.GetPlayers", nil, &ret)
	if err != nil {
		log.Fatalf("error getting players: %s", err)
	}
	return
}

func (c *Client) AddPlayer(player *Player) {
	err := c.client.Call("Group.AddPlayer", player, nil)
	if err != nil {
		log.Fatalf("error adding player: %s", err)
	}
	return
}
