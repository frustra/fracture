package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/frustra/fracture/chunk"
	"github.com/frustra/fracture/edge"
	"github.com/frustra/fracture/network"
)

func main() {
	var (
		// Cluster flags.
		existing = flag.String("join", "localhost:7946", "address of any node in the cluster")
		node     = flag.String("node", ":7946", "node address within the cluster")

		// Edge server flags.
		addr    = flag.String("addr", ":25565", "address to bind")
		players = flag.Int("size", 16, "max players")

		// Chunk server flags.
		// x = flag.Int64("x", 0, "x offset")
		// z = flag.Int64("z", 0, "z offset")
	)

	flag.Parse()
	role := flag.Arg(0)

	cluster, err := network.CreateCluster(*node, *existing)
	if err != nil {
		log.Fatal("Error creating cluster: ", err)
	}

	var server network.Server

	switch role {
	case "edge":
		server = &edge.Server{Addr: *addr, MaxPlayers: *players, Cluster: cluster}

	case "chunk":
		server = &chunk.Server{Addr: *addr, Cluster: cluster}

	default:
		log.Fatal("Invalid role: ", role)
	}

	cluster.NodeType = server.NodeType()
	cluster.NodeAddr = ":" + strconv.Itoa(server.NodePort())

	if err := cluster.Join(); err != nil {
		log.Fatal("Failed to join cluster: ", err)
	}

	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	go func() {
		for _ = range interrupts {
			cluster.Part()
			os.Exit(1)
		}
	}()

	log.Fatal(server.Serve())
}
