package main

import (
	glog "log"
	"sync"

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
	return s.firstIndex, nil
}

// LastIndex implements the LogStore interface.
func (s *logStore) LastIndex() (uint64, error) {
	return s.lastIndex, nil
}

// GetLog implements the LogStore interface.
func (s *logStore) GetLog(index uint64, log *raft.Log) error {
	glog.Printf("GetLog %d", index)
	return nil
}

// StoreLog implements the LogStore interface.
func (s *logStore) StoreLog(log *raft.Log) error {
	return s.StoreLogs([]*raft.Log{log})
}

// StoreLogs implements the LogStore interface.
func (s *logStore) StoreLogs(logs []*raft.Log) error {
	s.Lock()
	defer s.Unlock()

	for _, l := range logs {
		glog.Printf("StoreLogs %+v", l)

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
	glog.Printf("DeleteRange %d-%d", min, max)

	s.Lock()
	defer s.Unlock()

	for i := min; i <= max; i++ {

	}

	s.firstIndex = max + 1
	return nil
}

// Set implements the StableStore interface.
func (s *logStore) Get(key []byte) ([]byte, error) {
	glog.Printf("Get %s", string(key))
	return nil, nil
}

// GetUint64 implements the StableStore interface.
func (s *logStore) GetUint64(key []byte) (uint64, error) {
	glog.Printf("GetUint64 %s", string(key))
	// e,g.
	// GetUint64 CurrentTerm

	return 0, nil
}

// Set implements the StableStore interface.
func (s *logStore) Set(key, val []byte) error {
	glog.Printf("Set %s:%s", string(key), string(val))
	// e,g.
	// Set LastVoteCand:10.1.1.1:10114

	return nil
}

// SetUint64 implements the StableStore interface.
func (s *logStore) SetUint64(key []byte, val uint64) error {
	glog.Printf("SetUint64 %s:%d", string(key), val)
	// e,g.
	// SetUint64 CurrentTerm:1
	// SetUint64 LastVoteTerm:1

	return nil
}
