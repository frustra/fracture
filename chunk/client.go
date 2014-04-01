package chunk

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

func (c *Client) GetChunk(y uint) (ret *Chunk) {
	err := c.client.Call("Column.GetChunk", y, &ret)
	if err != nil {
		log.Fatalf("error getting chunk at %d: %s", y, err)
	}
	return
}

func (c *Client) GetChunks() (ret *[16]*Chunk) {
	err := c.client.Call("Column.GetChunks", nil, &ret)
	if err != nil {
		log.Fatalf("error getting chunks: %s", err)
	}
	return
}

func (c *Client) SetBlock(x, y, z, id uint8) {
	args := &Block{X: x, Y: y, Z: z, Id: id}
	err := c.client.Call("Column.SetBlock", args, nil)
	if err != nil {
		log.Fatalf("error setting chunk at %d, %d, %d: %s", x, y, z, err)
	}
}
