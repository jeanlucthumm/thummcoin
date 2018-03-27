package main

import (
	"github.com/jeanlucthumm/thummcoin/node"
	"net"
	"flag"
)

func main() {
	var seedMode = flag.Bool("seed", false, "enable seeding mode")
	flag.Parse()

	n := node.NewNode()
	addr, _ := net.ResolveTCPAddr("tcp", ":8080")
	n.Start(addr)

	if !*seedMode {
		n.Discover()
	}
}
