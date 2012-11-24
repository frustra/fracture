package edge

import (
	"log"
	"fracture/chunk"
)

func Serve() {
	c := new(chunk.Client)
	c.Connect("localhost:12444")
	c.SetBlock(0, 0, 0, 2)
	ch := c.GetChunk(0)
	log.Printf("%d should equal 2", ch.Data[0][0][0])
	log.Println("Starting up edge server")
}
