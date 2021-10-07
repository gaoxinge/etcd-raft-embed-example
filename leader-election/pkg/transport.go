package pkg

import (
	"context"
	"log"

	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
)

type Transport struct {
	id   uint64
	m    map[uint64]chan raftpb.Message
	node raft.Node
}

func NewTransport(id uint64, m map[uint64]chan raftpb.Message, node raft.Node) *Transport {
	return &Transport{
		id:   id,
		m:    m,
		node: node,
	}
}

func (t *Transport) Run() {
	for {
		select {
		case message := <-t.m[t.id]:
			log.Printf("[INFO] %d recv message %v", t.id, message)
			if err := t.node.Step(context.TODO(), message); err != nil {
				log.Printf("[ERROR] %d step message %v with error %v", t.id, message, err)
			}
		}
	}
}

func (t *Transport) Send(messages []raftpb.Message) {
	for _, message := range messages {
		log.Printf("[INFO] %d send message %v", t.id, message)
		t.m[message.To] <- message
	}
}
