package main

import (
	glog "log"
	"sync"
	"time"

	debugger "github.com/funkygao/golib/debug"
	"github.com/hashicorp/raft"
)

var (
	_ raft.LogStore    = &logStore{}
	_ raft.StableStore = &logStore{}
)

// logStore implements raft.LogStore and raft.StableStore.
type logStore struct {
	sync.RWMutex

	firstIndex, lastIndex uint64
}

func NewStore() *logStore {
	return &logStore{}
}

// FirstIndex implements the LogStore interface.
func (s *logStore) FirstIndex() (uint64, error) {
	glog.Printf("FirstIndex %+v", debugger.Callstack(2))

	return s.firstIndex, nil
}

// LastIndex implements the LogStore interface.
func (s *logStore) LastIndex() (uint64, error) {
	glog.Printf("LastIndex %+v", debugger.Callstack(2))

	return s.lastIndex, nil
}

// GetLog implements the LogStore interface.
func (s *logStore) GetLog(index uint64, log *raft.Log) error {
	glog.Printf("GetLog %d %+v", index, debugger.Callstack(2))
	time.Sleep(time.Millisecond * 100)
	return nil
}

// StoreLog implements the LogStore interface.
func (s *logStore) StoreLog(log *raft.Log) error {
	glog.Printf("StoreLog %+v", debugger.Callstack(2))
	return s.StoreLogs([]*raft.Log{log})
}

// StoreLogs implements the LogStore interface.
func (s *logStore) StoreLogs(logs []*raft.Log) error {
	glog.Println(len(logs), debugger.Callstack(2))

	s.Lock()
	defer s.Unlock()

	for _, l := range logs {
		glog.Printf("StoreLogs {idx:%d, term:%d, type:%d, data:%s}", l.Index, l.Term, l.Type, string(l.Data))

		if s.firstIndex == 0 {
			s.firstIndex = l.Index
		}
		if l.Index > s.lastIndex {
			s.lastIndex = l.Index
		}
	}

	return nil
}

// DeleteRange implements the LogStore interface.
func (s *logStore) DeleteRange(min, max uint64) error {
	glog.Printf("DeleteRange %d-%d %+v", min, max, debugger.Callstack(2))

	s.Lock()
	defer s.Unlock()

	for i := min; i <= max; i++ {

	}

	s.firstIndex = max + 1
	return nil
}

// Set implements the StableStore interface.
func (s *logStore) Get(key []byte) ([]byte, error) {
	glog.Printf("Get %s %+v", string(key), debugger.Callstack(2))
	return nil, nil
}

// GetUint64 implements the StableStore interface.
func (s *logStore) GetUint64(key []byte) (uint64, error) {
	glog.Printf("GetUint64 %s %+v", string(key), debugger.Callstack(2))
	// e,g.
	// GetUint64 CurrentTerm

	return 0, nil
}

// Set implements the StableStore interface.
func (s *logStore) Set(key, val []byte) error {
	glog.Printf("Set %s:%s %+v", string(key), string(val), debugger.Callstack(2))
	// e,g.
	// Set LastVoteCand:10.1.1.1:10114

	return nil
}

// SetUint64 implements the StableStore interface.
func (s *logStore) SetUint64(key []byte, val uint64) error {
	glog.Printf("SetUint64 %s=%d %+v", string(key), val, debugger.Callstack(2))
	// e,g.
	// SetUint64 CurrentTerm:1
	// SetUint64 LastVoteTerm:1

	return nil
}
