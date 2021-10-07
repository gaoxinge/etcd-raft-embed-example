package pkg

import (
	"log"
	"time"

	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
)

type Wrapper struct {
	id        uint64
	peers     []raft.Peer
	node      raft.Node
	storage   *raft.MemoryStorage
	transport *Transport
	ticker    *time.Ticker
}

func NewWrapper(id uint64, peers []raft.Peer, m map[uint64]chan raftpb.Message) *Wrapper {
	storage := raft.NewMemoryStorage()
	config := &raft.Config{
		ID:              id,
		ElectionTick:    10,
		HeartbeatTick:   1,
		Storage:         storage,
		MaxSizePerMsg:   4096,
		MaxInflightMsgs: 256,
	}
	node := raft.StartNode(config, peers)
	transport := NewTransport(id, m, node)
	ticker := time.NewTicker(5 * time.Second)
	return &Wrapper{
		id:        id,
		peers:     peers,
		node:      node,
		storage:   storage,
		transport: transport,
		ticker:    ticker,
	}
}

func (w *Wrapper) RunTransport() {
	w.transport.Run()
}

func (w *Wrapper) Run() {
	for {
		select {
		case <-w.ticker.C:
			log.Printf("[INFO] %d tick", w.id)
			w.node.Tick()
		case rd := <-w.node.Ready():
			log.Printf("[INFO] %d handle ready %v", w.id, rd)
			w.storage.Append(rd.Entries)
			w.transport.Send(rd.Messages)
			w.node.Advance()
		}
	}
}
