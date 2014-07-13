package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/frustra/fracture/chunk"
	"github.com/frustra/fracture/edge"
	"github.com/frustra/fracture/entity"
	"github.com/frustra/fracture/network"
	"github.com/frustra/fracture/world"
)

func main() {
	var (
		// Cluster flags.
		existing = flag.String("join", "localhost:7946", "address of any node in the cluster")
		node     = flag.String("node", ":7946", "node address within the cluster")

		// Edge server flags.
		addr   = flag.String("addr", ":25565", "address to bind")
		offset = flag.Int64("offset", 0, "id range offset")
		size   = flag.Int("size", 16, "server size")

		// Chunk server flags.
		px = flag.Int64("x", 0, "x offset")
		pz = flag.Int64("z", 0, "z offset")
	)

	log.SetFlags(log.Lmicroseconds)
	flag.Parse()
	role := flag.Arg(0)

	cluster, err := network.CreateCluster(*node, *existing)
	if err != nil {
		log.Fatal("Error creating cluster: ", err)
	}

	var server network.Server
	meta := cluster.LocalNodeMeta

	switch role {
	case "edge":
		meta.Type = network.EdgeType
		server = &edge.Server{Addr: *addr, Cluster: cluster, Size: *size, Offset: *offset}

	case "entity":
		meta.Type = network.EntityType
		server = &entity.Server{Addr: *addr, Cluster: cluster, Size: *size}

	case "chunk":
		meta.Type = network.ChunkType
		x, z := *px, *pz

		chunkServer := &chunk.Server{
			Addr:    *addr,
			Cluster: cluster,
			OffsetX: x * world.ChunkWidthPerNode,
			OffsetZ: z * world.ChunkWidthPerNode,
		}

		meta.X = &chunkServer.OffsetX
		meta.Z = &chunkServer.OffsetZ
		server = chunkServer

	default:
		log.Fatal("Invalid role: ", role)
	}

	log.SetPrefix(fmt.Sprintf("[%7s %-7s] ", role, *addr))
	meta.Addr = ":" + strconv.Itoa(server.NodePort())

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
