package main

import (
	"log"
	"flag"
	"github.com/frustra/fracture/chunk"
	"github.com/frustra/fracture/edge"
	"github.com/frustra/fracture/master"
	"github.com/frustra/fracture/player"
)

var (
	role = flag.String("role", "master", "server role")
	addr = flag.String("addr", "", "address to bind to")

	maxPlayers = flag.Int("max", 16, "max players")

	x = flag.Int64("x", 0, "x offset")
	z = flag.Int64("z", 0, "z offset")
)

func main() {
	flag.Parse()
	addr := *addr

	switch *role {
	case "master":
		if addr == "" {
			addr = ":25566"
		}
		master.Serve(addr)

	case "chunk":
		if addr == "" {
			log.Fatal("Must specify an addr")
		}
		chunk.Serve(*x, *z, addr)

	case "edge":
		edge.Serve()

	case "player":
		if addr == "" {
			addr = ":12444"
		}
		player.Serve(*maxPlayers, addr)

	default:
		log.Fatal("Invalid role", role)
	}
}
