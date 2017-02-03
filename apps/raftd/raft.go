package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/funkygao/gafka/ctx"
	"github.com/hashicorp/raft"
)

func MakeRaft() *raft.Raft {
	baseDir := "/tmp/raftd"
	stableStore := raft.NewInmemStore()
	snapshotStore, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stdout)
	if err != nil {
		panic(err)
	}

	ip, _ := ctx.LocalIP()

	trans, err := raft.NewTCPTransport(fmt.Sprintf("%s:10114", ip.String()), nil, 2, time.Second, os.Stdout)
	if err != nil {
		panic(err)
	}
	peerStore := raft.NewJSONPeers(baseDir, trans)
	log.Printf("Starting node at %v", trans.LocalAddr())

	conf := raft.DefaultConfig()
	conf.SnapshotInterval = time.Minute
	conf.EnableSingleNode = false
	conf.Logger = log.New(os.Stdout, "", log.LstdFlags)
	node, err := raft.NewRaft(conf, &MockFSM{}, stableStore, stableStore, snapshotStore, peerStore, trans)
	if err != nil {
		panic(err)
	}

	return node
}
