package main

import (
	"github.com/jeanlucthumm/thummcoin/node"
	"net"
	"flag"
	"log"
)

func main() {
	var seedMode = flag.Bool("seed", false, "enable seeding mode")
	flag.Parse()

	n := node.NewNode()
	addr, _ := net.ResolveTCPAddr("tcp", ":8080")

	log.Println()

	if *seedMode {
		n.StartSeed(addr)
	} else {
		n.Start(addr)
	}
}
