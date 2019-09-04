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

	n := node.NewNode(*seedMode)
	addr, _ := net.ResolveTCPAddr("tcp", ":8080")

	log.Println()

	if err := n.Start(addr); err != nil {
		log.Printf("Failed to start node: %s\n", err)
	}
}
