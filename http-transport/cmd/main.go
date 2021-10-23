package main

import (
	"github.com/gaoxinge/etcd-raft-embed-example/http-transport/pkg"
)

func main() {
	peers := map[uint64]string{
		0x01: "http-transport://127.0.0.1:20001",
		0x02: "http-transport://127.0.0.1:20002",
		0x03: "http-transport://127.0.0.1:20003",
	}

	wrapper1 := pkg.NewWrapper(0x01, peers)
	go wrapper1.RunTransport()
	go wrapper1.Run()

	wrapper2 := pkg.NewWrapper(0x02, peers)
	go wrapper2.RunTransport()
	go wrapper2.Run()

	wrapper3 := pkg.NewWrapper(0x03, peers)
	go wrapper3.RunTransport()
	go wrapper3.Run()

	select {}
}
