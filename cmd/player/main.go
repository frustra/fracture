package main

import (
	"flag"
	"github.com/frustra/fracture/player"
)

var maxPlayers = flag.Int("max", 16, "max players")
var addr = flag.String("addr", ":12444", "address to bind to")

func main() {
	flag.Parse()
	player.Serve(*maxPlayers, *addr)
}
