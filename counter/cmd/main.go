package main

import (
	"log"

	"github.com/gaoxinge/etcd-raft-embed-example/counter/pkg"
)

func main() {
	var err error

	peers := map[uint64]string{
		0x01: "http://127.0.0.1:20001",
		0x02: "http://127.0.0.1:20002",
		0x03: "http://127.0.0.1:20003",
	}

	proposeC1 := make(chan struct{})
	commitC1 := make(chan uint64)
	_, err = pkg.NewWrapper(0x01, peers, proposeC1, commitC1)
	if err != nil {
		log.Printf("[ERROR] %d init with error %v", 0x01, err)
		return
	}
	counter1 := pkg.NewCounter(proposeC1, commitC1)
	pkg.NewServer("http://127.0.0.1:30001", counter1)

	proposeC2 := make(chan struct{})
	commitC2 := make(chan uint64)
	_, err = pkg.NewWrapper(0x02, peers, proposeC2, commitC2)
	if err != nil {
		log.Printf("[ERROR] %d init with error %v", 0x02, err)
		return
	}
	counter2 := pkg.NewCounter(proposeC2, commitC2)
	pkg.NewServer("http://127.0.0.1:30002", counter2)

	proposeC3 := make(chan struct{})
	commitC3 := make(chan uint64)
	_, err = pkg.NewWrapper(0x03, peers, proposeC3, commitC3)
	if err != nil {
		log.Printf("[ERROR] %d init with error %v", 0x03, err)
		return
	}
	counter3 := pkg.NewCounter(proposeC3, commitC3)
	pkg.NewServer("http://127.0.0.1:30003", counter3)

	select {}
}
