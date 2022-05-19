package pkg

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"sync"
)

const (
	PUT = 0
	GET = 1
)

type status struct {
	s  bool
	v  string
	ok bool
}

type entry struct {
	OP    int
	K     string
	V     string
	Id    uint64
	Index uint64
}

func (e entry) String() string {
	return fmt.Sprintf("entry{OP: %d, K: %s, V: %s, Id: %d, Index: %d}", e.OP, e.K, e.V, e.Id, e.Index)
}

type KV struct {
	id          uint64
	index       uint64
	storage     *Storage
	kv          map[string]string
	waitMap     map[uint64]chan status
	waitMapLock sync.Mutex
	proposeC    chan<- string
	commitC     <-chan []string
	mu          sync.Mutex
}

func NewKV(id uint64, proposeC chan<- string, commitC <-chan []string) (*KV, error) {
	storagePath := fmt.Sprintf("linear-kv/storage/%d", id)
	storage, err := NewStorage(storagePath)
	if err != nil {
		return nil, err
	}
	err = storage.Create(0)
	if err != nil {
		return nil, err
	}
	index, err := storage.Get()
	if err != nil {
		return nil, err
	}

	kv := &KV{
		id:       id,
		index:    index,
		storage:  storage,
		kv:       make(map[string]string),
		waitMap:  make(map[uint64]chan status),
		proposeC: proposeC,
		commitC:  commitC,
	}
	go kv.readCommits()
	return kv, nil
}

func (kv *KV) IncIndex() error {
	err := kv.storage.Put(kv.index + 1)
	if err != nil {
		return err
	}
	kv.index += 1
	return nil
}

func (kv *KV) Put(k string, v string) bool {
	waitCh := make(chan status)

	kv.mu.Lock()
	if err := kv.IncIndex(); err != nil {
		log.Fatalf("[ERROR] %d put with error %v", kv.id, err)
	}

	kv.waitMapLock.Lock()
	kv.waitMap[kv.index] = waitCh
	kv.waitMapLock.Unlock()

	var buf bytes.Buffer
	e := entry{
		OP:    PUT,
		K:     k,
		V:     v,
		Id:    kv.id,
		Index: kv.index,
	}
	if err := gob.NewEncoder(&buf).Encode(e); err != nil {
		log.Fatalf("[ERROR] %d put with error %v", kv.id, err)
	}

	log.Printf("[INFO] %d put %s", kv.id, e)

	kv.proposeC <- buf.String()
	kv.mu.Unlock()

	s := <-waitCh
	return s.s
}

func (kv *KV) Get(k string) (bool, string, bool) {
	waitCh := make(chan status)

	kv.mu.Lock()
	if err := kv.IncIndex(); err != nil {
		log.Fatalf("[ERROR] %d get with error %v", kv.id, err)
	}

	kv.waitMapLock.Lock()
	kv.waitMap[kv.index] = waitCh
	kv.waitMapLock.Unlock()

	var buf bytes.Buffer
	e := entry{
		OP:    GET,
		K:     k,
		V:     "",
		Id:    kv.id,
		Index: kv.index,
	}
	if err := gob.NewEncoder(&buf).Encode(e); err != nil {
		log.Fatalf("[ERROR] %d get with error %v", kv.id, err)
	}

	log.Printf("[INFO] %d get %s", kv.id, e)

	kv.proposeC <- buf.String()
	kv.mu.Unlock()

	s := <-waitCh
	return s.s, s.v, s.ok
}

func (kv *KV) readCommits() {
	for datas := range kv.commitC {
		for _, data := range datas {
			var e entry
			if err := gob.NewDecoder(bytes.NewBufferString(data)).Decode(&e); err != nil {
				log.Fatalf("[ERROR] %d read commits with error %v", kv.id, err)
			}

			log.Printf("[INFO] %d handle %s", kv.id, e)

			var v string
			var ok bool
			if e.OP == PUT {
				kv.kv[e.K] = e.V
				v = ""
				ok = false
			} else {
				v, ok = kv.kv[e.K]
			}

			if e.Id == kv.id {
				kv.waitMapLock.Lock()
				waitCh, exist := kv.waitMap[e.Index]
				if exist {
					waitCh <- status{
						s:  true,
						v:  v,
						ok: ok,
					}
					close(waitCh)
					delete(kv.waitMap, e.Index)
				}
				kv.waitMapLock.Unlock()

				kv.waitMapLock.Lock()
				for i := range kv.waitMap {
					if i < e.Index {
						waitCh := kv.waitMap[i]
						waitCh <- status{
							s:  false,
							v:  "",
							ok: false,
						}
						close(waitCh)
						delete(kv.waitMap, i)
					}
				}
				kv.waitMapLock.Unlock()
			}
		}
	}
}
