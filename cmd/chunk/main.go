package main

import (
	"fracture/chunk"
	"flag"
)

var x = flag.Int64("x", 0, "x offset")
var z = flag.Int64("z", 0, "z offset")
var addr = flag.String("addr", ":12444", "address to bind to")

func main() {
	flag.Parse()
	chunk.Serve(*x, *z, *addr)
}
