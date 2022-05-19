package main

import (
	"log"

	"github.com/gaoxinge/etcd-raft-embed-example/linear-kv/pkg"
)

func main() {
	var err error

	peers := map[uint64]string{
		0x01: "http://127.0.0.1:20001",
		0x02: "http://127.0.0.1:20002",
		0x03: "http://127.0.0.1:20003",
	}

	proposeC1 := make(chan string)
	commitC1 := make(chan []string)
	_, err = pkg.NewWrapper(0x01, peers, proposeC1, commitC1)
	if err != nil {
		log.Printf("[ERROR] %d init with error %v", 0x01, err)
		return
	}
	kv1, err := pkg.NewKV(0x01, proposeC1, commitC1)
	if err != nil {
		log.Printf("[ERROR] %d init with error %v", 0x01, err)
		return
	}
	pkg.NewServer("http://127.0.0.1:30001", kv1)

	proposeC2 := make(chan string)
	commitC2 := make(chan []string)
	_, err = pkg.NewWrapper(0x02, peers, proposeC2, commitC2)
	if err != nil {
		log.Printf("[ERROR] %d init with error %v", 0x02, err)
		return
	}
	kv2, err := pkg.NewKV(0x02, proposeC2, commitC2)
	if err != nil {
		log.Printf("[ERROR] %d init with error %v", 0x02, err)
		return
	}
	pkg.NewServer("http://127.0.0.1:30002", kv2)

	proposeC3 := make(chan string)
	commitC3 := make(chan []string)
	_, err = pkg.NewWrapper(0x03, peers, proposeC3, commitC3)
	if err != nil {
		log.Printf("[ERROR] %d init with error %v", 0x03, err)
		return
	}
	kv3, err := pkg.NewKV(0x03, proposeC3, commitC3)
	if err != nil {
		log.Printf("[ERROR] %d init with error %v", 0x03, err)
		return
	}
	pkg.NewServer("http://127.0.0.1:30003", kv3)

	select {}
}
