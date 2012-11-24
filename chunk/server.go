package chunk

import(
	"net"
	"net/http"
	"net/rpc"
	"log"
	"time"
)

var col *Column

func Serve(x int64, z int64, addr string) {
	col = new(Column)
	for i := 0; i < 16; i++ {
		col.chunks[i] = new(Chunk)
	}
	rpc.Register(col)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting up chunk server on %s for (%d, %d)", addr, x, z)

	go func() {
		for i := 0; i < 16; i++ {
			col.chunks[i].Tick()
		}
		time.Sleep(50 * time.Millisecond) // 20 ticks/sec
	}()

	http.Serve(l, nil)
}
