package main

import (
	"flag"
	"log"

	"github.com/frustra/fracture/edge"
)

func main() {
	var (
		role = flag.String("role", "master", "server role")

		addr       = flag.String("addr", "", "address to bind to")
		maxPlayers = flag.Int("size", 16, "max players")

		// x = flag.Int64("x", 0, "x offset")
		// z = flag.Int64("z", 0, "z offset")
	)

	flag.Parse()

	switch *role {
	case "edge":
		edge.Serve(*addr, *maxPlayers)

	default:
		log.Fatal("Invalid role: ", *role)
	}
}
