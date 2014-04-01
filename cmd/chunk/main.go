package main

import (
	"flag"
	"github.com/frustra/fracture/chunk"
)

var x = flag.Int64("x", 0, "x offset")
var z = flag.Int64("z", 0, "z offset")
var addr = flag.String("addr", "127.0.0.1:0", "address to bind to")

func main() {
	flag.Parse()
	chunk.Serve(*x, *z, *addr)
}
