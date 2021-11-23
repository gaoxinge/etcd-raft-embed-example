package pkg

import (
	"bytes"
	"encoding/gob"
	"log"
	"sync"
)

type KV struct {
	id       uint64
	kv       map[string]string
	mu       sync.RWMutex
	proposeC chan<- string
	commitC  <-chan []string
}

func NewKV(id uint64, proposeC chan<- string, commitC <-chan []string) *KV {
	kv := &KV{
		id:       id,
		kv:       make(map[string]string),
		proposeC: proposeC,
		commitC:  commitC,
	}
	go kv.readCommits()
	return kv
}

func (kv *KV) Put(k string, v string) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(
		struct{
			K string
			V string
		}{
			k,
			v,
		}); err != nil {
		log.Fatalf("[ERROR] %d put with error %v", kv.id, err)
	}
	kv.proposeC <- buf.String()
}

func (kv *KV) Get(k string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	v, ok := kv.kv[k]
	return v, ok
}

func (kv *KV) readCommits() {
	for datas := range kv.commitC {
		for _, data := range datas {
			var p struct{
				K string
				V string
			}
			if err := gob.NewDecoder(bytes.NewBufferString(data)).Decode(&p); err != nil {
				log.Fatalf("[ERROR] %d read commits with error %v", kv.id, err)
			}
			kv.mu.Lock()
			kv.kv[p.K] = p.V
			kv.mu.Unlock()
		}
	}
}
