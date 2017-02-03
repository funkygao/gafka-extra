package main

import (
	"io"
	"log"
	"sync"

	"github.com/hashicorp/go-msgpack/codec"
	"github.com/hashicorp/raft"
)

var (
	_ raft.FSM         = &MockFSM{}
	_ raft.FSMSnapshot = &MockSnapshot{}
)

type MockFSM struct {
	sync.Mutex
	logs [][]byte
}

func (fsm *MockFSM) Apply(l *raft.Log) interface{} {
	if debug {
		log.Printf("apply %+v", l)
	}

	fsm.Lock()
	defer fsm.Unlock()

	fsm.logs = append(fsm.logs, l.Data)
	return len(fsm.logs)
}

func (fsm *MockFSM) Snapshot() (raft.FSMSnapshot, error) {
	if debug {
		log.Println("snapshot")
	}

	fsm.Lock()
	defer fsm.Unlock()

	return &MockSnapshot{logs: fsm.logs, maxIdx: len(fsm.logs)}, nil
}

func (fsm *MockFSM) Restore(rc io.ReadCloser) error {
	if debug {
		log.Println("restore")
	}

	fsm.Lock()
	defer fsm.Unlock()
	defer rc.Close()

	dec := codec.NewDecoder(rc, &codec.MsgpackHandle{})
	fsm.logs = nil
	return dec.Decode(&fsm.logs)
}

type MockSnapshot struct {
	logs   [][]byte
	maxIdx int
}

func (snap *MockSnapshot) Persist(sink raft.SnapshotSink) error {
	if debug {
		log.Printf("persist %#v", sink)
	}

	enc := codec.NewEncoder(sink, &codec.MsgpackHandle{})
	if err := enc.Encode(snap.logs[:snap.maxIdx]); err != nil {
		sink.Cancel()
		return err
	}

	sink.Close()
	return nil
}

func (snap *MockSnapshot) Release() {
	if debug {
		log.Printf("release")
	}
}
