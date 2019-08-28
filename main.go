package main

import (
	"flag"
	"github.com/jeanlucthumm/thummcoin/node"
	"log"
	"net"
)

func main() {
	var seedMode = flag.Bool("seed", false, "enable seeding mode")
	flag.Parse()

	n := node.NewNode()
	addr, _ := net.ResolveTCPAddr("tcp", ":8080")

	log.Println()

	if *seedMode {
		if err := n.StartSeed(addr); err != nil {

		}
	} else {
		if err := n.Start(addr); err != nil {
			return
		}
	}
}
