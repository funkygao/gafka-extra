package main

import (
	"log"
	"os"
	"time"

	"github.com/hashicorp/raft"
)

func MakeRaft() *raft.Raft {
	baseDir := "/tmp/raftd"
	stableStore := raft.NewInmemStore()
	snapshotStore, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stdout)
	if err != nil {
		panic(err)
	}

	trans, err := raft.NewTCPTransport(":10114", nil, 2, time.Second, os.Stdout)
	if err != nil {
		panic(err)
	}
	peerStore := raft.NewJSONPeers(baseDir, trans)
	log.Printf("Starting node at %v", trans.LocalAddr())

	conf := raft.DefaultConfig()
	conf.Logger = log.New(os.Stdout, "", log.LstdFlags)
	node, err := raft.NewRaft(conf, &MockFSM{}, stableStore, stableStore, snapshotStore, peerStore, trans)
	if err != nil {
		panic(err)
	}

	return node
}
