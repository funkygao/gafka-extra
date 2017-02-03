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
	future := node.SetPeers(strings.Split(peer, ","))
	if err := future.Error(); err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(time.Second * 30)
		defer ticker.Stop()

		for range ticker.C {
			log.Printf("[%s] state:%+v stats:%+v", node, node.State(), node.Stats())
		}
	}()

	defer node.Shutdown()

	if writer {
		log.SetOutput(os.Stdout)
		stress.Flags.Round = 5
		stress.Flags.Tick = 5
		stress.RunStress(benchAppend)
	} else {
		select {}
	}
}

func benchAppend(seq int) {
	cmd := []byte("hello world")
	for i := 0; i < 1000; i++ {
		future := node.Apply(cmd, time.Second)
		if future.Error() == nil {
			stress.IncCounter("ok", 1)
		} else {
			stress.IncCounter("no", 1)
		}
	}
}
