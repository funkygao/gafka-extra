package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/funkygao/golib/stress"
	"github.com/hashicorp/raft"
)

var (
	debug  bool
	peer   string
	writer bool

	node *raft.Raft
)

func init() {
	flag.BoolVar(&debug, "d", false, "debug")
	flag.StringVar(&peer, "peers", "", "peers comma separated host:port")
	flag.BoolVar(&writer, "w", false, "this raft writes data")
	flag.Parse()

	if len(peer) < 1 {
		panic("peers is required")
	}
}

func main() {
	node = MakeRaft()
	log.Println("raft made, debug=%v peers=%s w=%v", debug, peer, writer)
	future := node.SetPeers(strings.Split(peer, ","))
	if err := future.Error(); err != nil {
		panic(err)
	}

	log.Println("raft set peers done")

	go func() {
		ticker := time.NewTicker(time.Second * 30)
		defer ticker.Stop()

		for range ticker.C {
			log.Printf("[%s] stats:%+v", node, node.Stats())
		}
	}()

	defer node.Shutdown()

	log.Println("enter main loop")

	if writer {
		log.SetOutput(os.Stdout)
		stress.Flags.Round = 5
		stress.Flags.Tick = 5
		stress.RunStress(benchAppend)
	} else {
		select {}
	}
}
