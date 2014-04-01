package main

import (
	"flag"
	"github.com/frustra/fracture/master"
)

var addr = flag.String("addr", "127.0.0.1:25566", "address to bind to")

func main() {
	flag.Parse()
	master.Serve(*addr)
}
