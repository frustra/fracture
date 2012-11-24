package player

import(
	"net"
	"net/http"
	"net/rpc"
	"log"
//	"time"
)

func Serve(maxPlayers int, addr string) {
	group := new(Group)
	group.maxPlayers = maxPlayers
	group.players = make([]*Player, 0, maxPlayers)
	rpc.Register(group)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting up player server on %s", addr)

	/*go func() {
		time.Sleep(50 * time.Millisecond) // 20 ticks/sec
	}()*/

	http.Serve(l, nil)
}
