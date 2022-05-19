package pkg

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	stats "go.etcd.io/etcd/server/v3/etcdserver/api/v2stats"
	"go.etcd.io/etcd/server/v3/wal"
	"go.etcd.io/etcd/server/v3/wal/walpb"
	"go.uber.org/zap"
)

type Wrapper struct {
	id    uint64
	peers map[uint64]string

	node      raft.Node
	wal       *wal.WAL
	storage   *raft.MemoryStorage
	transport *rafthttp.Transport
	ticker    *time.Ticker

	proposeC     <-chan string
	commitC      chan<- []string
	appliedIndex uint64
}

func NewWrapper(id uint64, peers map[uint64]string, proposeC <-chan string, commitC chan<- []string) (*Wrapper, error) {
	wrapper := &Wrapper{
		id:    id,
		peers: peers,
	}

	npeers := make([]raft.Peer, 0, len(peers))
	for i := range peers {
		npeers = append(npeers, raft.Peer{ID: i})
	}

	storage := raft.NewMemoryStorage()
	waldir := fmt.Sprintf("linear-kv/wal/%d", id)
	oldwal := wal.Exist(waldir)
	if !wal.Exist(waldir) {
		err := os.Mkdir(waldir, 0705)
		if err != nil {
			return nil, err
		}
		wal0, err := wal.Create(zap.NewExample(), waldir, nil)
		if err != nil {
			return nil, err
		}
		wal0.Close()
	}
	wal0, err := wal.Open(zap.NewExample(), waldir, walpb.Snapshot{})
	if err != nil {
		return nil, err
	}
	_, hardState, entries, err := wal0.ReadAll()
	if err != nil {
		return nil, err
	}
	storage.SetHardState(hardState)
	storage.Append(entries)

	config := &raft.Config{
		ID:              id,
		ElectionTick:    10,
		HeartbeatTick:   1,
		Storage:         storage,
		MaxSizePerMsg:   4096,
		MaxInflightMsgs: 256,
	}
	if oldwal {
		wrapper.node = raft.RestartNode(config)
	} else {
		wrapper.node = raft.StartNode(config, npeers)
	}

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

	wrapper.wal = wal0
	wrapper.storage = storage
	wrapper.transport = transport
	wrapper.ticker = time.NewTicker(5 * time.Second)
	wrapper.proposeC = proposeC
	wrapper.commitC = commitC

	go wrapper.RunPropose()
	go wrapper.RunTransport()
	go wrapper.Run()

	return wrapper, nil
}

func (w *Wrapper) RunPropose() {
	for {
		select {
		case data := <-w.proposeC:
			w.node.Propose(context.TODO(), []byte(data))
		}
	}
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
			w.wal.Save(rd.HardState, rd.Entries)
			w.storage.Append(rd.Entries)
			w.transport.Send(rd.Messages)

			entries := rd.CommittedEntries
			if len(entries) > 0 {
				entries = entries[w.appliedIndex+1-entries[0].Index:]
				datas := make([]string, 0, len(entries))
				for _, entry := range entries {
					switch entry.Type {
					case raftpb.EntryNormal:
						if len(entry.Data) == 0 {
							break
						}
						datas = append(datas, string(entry.Data))
					case raftpb.EntryConfChange:
						var cc raftpb.ConfChange
						cc.Unmarshal(entry.Data)
						w.node.ApplyConfChange(cc)
						switch cc.Type {
						case raftpb.ConfChangeAddNode:
							if len(cc.Context) > 0 {
								w.transport.AddPeer(types.ID(cc.NodeID), []string{string(cc.Context)})
							}
						case raftpb.ConfChangeRemoveNode:
							w.transport.RemovePeer(types.ID(cc.NodeID))
						}
					}
				}
				w.commitC <- datas
				w.appliedIndex = entries[len(entries)-1].Index
			}

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
