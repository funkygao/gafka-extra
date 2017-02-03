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
}

func main() {
	node = MakeRaft()
	if len(peer) > 0 {
		for _, p := range strings.Split(peer, ",") {
			log.Printf("add peer: %s", p)

			future := node.AddPeer(p)
			if err := future.Error(); err != nil {
				panic(err)
			}
		}
	}

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer ticker.Stop()

		for range ticker.C {
			log.Printf("state:%+v stats:%+v", node.State(), node.Stats())
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
