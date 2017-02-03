package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/funkygao/gafka/ctx"
	"github.com/hashicorp/raft"
)

func MakeRaft(baseDir string) *raft.Raft {
	// create the log store and stable store
	logStore := raft.NewInmemStore()
	stableStore := logStore

	// create the snapshot store, which allows raft to truncate the log
	snapshotStore, err := raft.NewFileSnapshotStore(baseDir, 3, os.Stdout)
	if err != nil {
		panic(err)
	}

	// setup raft communication
	ip, _ := ctx.LocalIP()
	raftBindAddr := fmt.Sprintf("%s:10114", ip.String())
	advertiseAddr, err := net.ResolveTCPAddr("tcp", raftBindAddr)
	if err != nil {
		panic(err)
	}
	maxPool := 3
	trans, err := raft.NewTCPTransport(raftBindAddr, advertiseAddr, maxPool, 4*time.Second, os.Stdout)
	if err != nil {
		panic(err)
	}
	log.Printf("Starting node at %v", trans.LocalAddr())

	// create peer storage
	peerStore := raft.NewJSONPeers(baseDir, trans)

	// check for any existing peers
	peers, err := peerStore.Peers()
	if err != nil {
		panic(err)
	}
	log.Printf("peers: %+v", peers)

	// setup the config
	conf := raft.DefaultConfig()
	conf.SnapshotInterval = time.Minute
	conf.EnableSingleNode = false
	conf.Logger = log.New(os.Stdout, "", log.LstdFlags)

	// create the raft system
	node, err := raft.NewRaft(conf, NewFSM(), logStore, stableStore, snapshotStore, peerStore, trans)
	if err != nil {
		panic(err)
	}

	return node
}
