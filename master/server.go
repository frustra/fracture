package master

import(
	"net"
	"net/http"
	"net/rpc"
	"log"
//	"time"
)

// Variable naming is shit in this file. Ignore it for now.


type ChunkServer struct {
	X, Z int64
	Addr string
}

type EdgeServer struct {
	Addr string
	client *rpc.Client
}

func (e *EdgeServer) BroadcastAddedChunk(c *ChunkServer) {
	e.client.Call("ChunkNodes.AddChunkNode", c, nil)
}

func (e *EdgeServer) BroadcastRemovedChunk(c *ChunkServer) {
	e.client.Call("ChunkNodes.RemoveChunkNode", c, nil)
}


type Server struct {
	Chunks map[string]*ChunkServer
	Edges map[string]*EdgeServer
}

func (s *Server) AnnounceChunkServer(c *ChunkServer, unused *bool) error {
	for _, e := range s.Edges {
		go e.BroadcastAddedChunk(c)
	}
	s.Chunks[c.Addr] = c
	log.Printf("Chunk server %s connected", c.Addr)
	return nil
}

func (s *Server) DestroyChunkServer(c *ChunkServer, unused *bool) error {
	cc, exists := s.Chunks[c.Addr]
	if !exists {
		return nil
	}
	delete(s.Chunks, cc.Addr)
	for _, e := range s.Edges {
		go e.BroadcastRemovedChunk(cc)
	}
	log.Printf("Chunk server %s disconnected", c.Addr)
	return nil
}

func (s *Server) AnnounceEdgeServer(e *EdgeServer, unused *bool) error {
	s.Edges[e.Addr] = e
	log.Printf("Edge server %s connected", e.Addr)
	go func() {
		var err error
		e.client, err = rpc.DialHTTP("tcp", e.Addr)
		if err != nil {
			log.Printf("Failed to connect to edge server")
			return
		}
		for _, c := range s.Chunks {
			go e.BroadcastAddedChunk(c)
		}
	}()
	return nil
}

func (s *Server) DestroyEdgeServer(e *EdgeServer, unused *bool) error {
	edge, exists := s.Edges[e.Addr]
	if !exists {
		return nil
	}
	edge.client.Close()
	delete(s.Edges, edge.Addr)
	log.Printf("Edge server %s disconnected", edge.Addr)
	return nil
}

func Serve(addr string) {
	server := new(Server)
	server.Chunks = make(map[string]*ChunkServer)
	server.Edges = make(map[string]*EdgeServer)
	rpc.Register(server)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting up master server on %s", addr)

	http.Serve(l, nil)
}
