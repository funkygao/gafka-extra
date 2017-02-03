package main

import (
	"time"

	"github.com/funkygao/golib/stress"
)

func benchAppend(seq int) {
	cmd := []byte("hello world")
	for i := 0; i < benchmark; i++ {
		if !isLeader.Get() {
			return
		}

		future := node.Apply(cmd, time.Second)
		if future.Error() == nil {
			stress.IncCounter("ok", 1)
		} else {
			stress.IncCounter("no", 1)
		}
	}
}
