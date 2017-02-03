package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/funkygao/golib/stress"
	"github.com/funkygao/golib/sync2"
	"github.com/hashicorp/raft"
)

var (
	debug     bool
	peer      string
	baseDir   string
	benchmark int

	isLeader sync2.AtomicBool
	node     *raft.Raft
)

func init() {
	flag.BoolVar(&debug, "d", false, "debug")
	flag.StringVar(&peer, "peers", "", "peers comma separated host:port")
	flag.StringVar(&baseDir, "base", "/tmp/raftd", "raft log dir")
	flag.IntVar(&benchmark, "benchn", 10000, "benchmark loop")
	flag.Parse()

	if len(peer) < 1 {
		panic("peers is required")
	}
}

func main() {
	node = MakeRaft(baseDir)
	log.Printf("raft made, debug=%v peers=%s w=%v", debug, peer)
	future := node.SetPeers(strings.Split(peer, ","))
	if err := future.Error(); err != nil {
		panic(err)
	}

	log.Println("raft set peers done")

	go func() {
		ticker := time.NewTicker(time.Second * 30)
		defer ticker.Stop()

		for range ticker.C {
			log.Printf("[%s] leader:%s stats:%+v", node, node.Leader(), node.Stats())
		}
	}()

	defer node.Shutdown()

	// setup the stress config
	log.SetOutput(os.Stdout)
	stress.Flags.Round = 5
	stress.Flags.Tick = 5

	log.Println("enter main loop")
	for {
		select {
		case leader := <-node.LeaderCh():
			if leader {
				log.Println("become leader, will start benchmark...")
				isLeader.Set(true)
				go stress.RunStress(benchAppend)
			} else {
				log.Println("lost leader state, will follow master...")
				isLeader.Set(false)
				time.Sleep(time.Second * 5) // await stress done
			}
		}
	}

}
