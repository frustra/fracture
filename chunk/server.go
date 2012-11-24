package chunk

import(
	"net"
	"net/http"
	"net/rpc"
	"log"
)

var col *Column

func Serve() {
	col = new(Column)
	rpc.Register(col)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal(e)
	}
	log.Println("Starting up chunk server")
	http.Serve(l, nil)
}
