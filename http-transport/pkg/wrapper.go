package pkg

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	stats "go.etcd.io/etcd/server/v3/etcdserver/api/v2stats"
	"go.uber.org/zap"
)

type Wrapper struct {
	id        uint64
	peers     map[uint64]string
	node      raft.Node
	storage   *raft.MemoryStorage
	transport *rafthttp.Transport
	ticker    *time.Ticker
}

func NewWrapper(id uint64, peers map[uint64]string) *Wrapper {
	wrapper := &Wrapper{
		id:     id,
		peers:  peers,
		ticker: time.NewTicker(5 * time.Second),
	}

	npeers := make([]raft.Peer, 0, len(peers))
	for i := range peers {
		npeers = append(npeers, raft.Peer{ID: i})
	}
	storage := raft.NewMemoryStorage()
	config := &raft.Config{
		ID:              id,
		ElectionTick:    10,
		HeartbeatTick:   1,
		Storage:         storage,
		MaxSizePerMsg:   4096,
		MaxInflightMsgs: 256,
	}
	node := raft.StartNode(config, npeers)
	transport := &rafthttp.Transport{
		Logger:      zap.NewExample(),
		ID:          types.ID(id),
		ClusterID:   0x1000,
		Raft:        wrapper,
		ServerStats: stats.NewServerStats("", ""),
		LeaderStats: stats.NewLeaderStats(zap.NewExample(), strconv.Itoa(int(id))),
		ErrorC:      make(chan error),
	}
	transport.Start()
	for i, addr := range peers {
		if i != id {
			transport.AddPeer(types.ID(i), []string{addr})
		}
	}

	wrapper.node = node
	wrapper.storage = storage
	wrapper.transport = transport

	return wrapper
}

func (w *Wrapper) RunTransport() {
	url, _ := url.Parse(w.peers[w.id])
	http.ListenAndServe(url.Host, w.transport.Handler())
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

func (w *Wrapper) Process(ctx context.Context, m raftpb.Message) error {
	return w.node.Step(ctx, m)
}

func (w *Wrapper) IsIDRemoved(id uint64) bool {
	return false
}

func (w *Wrapper) ReportUnreachable(id uint64) {

}

func (w *Wrapper) ReportSnapshot(id uint64, status raft.SnapshotStatus) {

}
