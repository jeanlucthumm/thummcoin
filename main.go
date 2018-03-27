package main

import (
	"github.com/jeanlucthumm/thummcoin/node"
	"net"
)

func main() {
	n := node.NewNode()
	addr, _ := net.ResolveTCPAddr("tcp", ":8080")
	n.Start(addr)
	n.Discover()
}
