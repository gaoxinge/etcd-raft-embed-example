package main

import (
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"

	"github.com/gaoxinge/etcd-raft-embed-example/leader-election/pkg"
)

func main() {
	peers := []raft.Peer{
		{ID: 0x01},
		{ID: 0x02},
		{ID: 0x03},
	}
	m := map[uint64]chan raftpb.Message{
		0x01: make(chan raftpb.Message),
		0x02: make(chan raftpb.Message),
		0x03: make(chan raftpb.Message),
	}

	wrapper1 := pkg.NewWrapper(0x01, peers, m)
	go wrapper1.RunTransport()
	go wrapper1.Run()

	wrapper2 := pkg.NewWrapper(0x02, peers, m)
	go wrapper2.RunTransport()
	go wrapper2.Run()

	wrapper3 := pkg.NewWrapper(0x03, peers, m)
	go wrapper3.RunTransport()
	go wrapper3.Run()

	select {}
}
