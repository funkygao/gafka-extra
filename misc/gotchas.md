# Gotchas

# zk

- ZooKeeper can start to perform badly if there are many nodes with thousands of children.
- The ZooKeeper database is kept entirely in memory.
- If a ZNode gets too big it can be extremely difficult to clean. 

  getChildren() will fail on the node. At Netflix we had to create a special-purpose program that had a huge value for jute.maxbuffer in order to get the nodes and delete them.

- ZooKeeper can slow down considerably on startup if there are many large ZNodes
- ZooKeeper has a 1MB transport limitation. In practice this means that ZNodes must be relatively small. 

- ZooKeeper watches are single threaded.

  No other watchers can be processed while your watcher is running

- Lock

  Very low session timeouts should be considered risky.

  https://qnalist.com/questions/6134306/locking-leader-election-and-dealing-with-session-loss

    A acquires lock ok
    B acquires lok fails, and blocked waiting for lock to release
    A STW GC beyond session timeout
    ZK will delete A ephemeral lock znode
    B watcher fired and acquires the lock ok
    A GC finishes: both A and B think they are the lock owner

