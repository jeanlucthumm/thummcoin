package main

import (
	"flag"
	"github.com/jeanlucthumm/thummcoin/seeder"
	"github.com/jeanlucthumm/thummcoin/node"
	"strconv"
)

func main() {
	seed := flag.Bool("seed", false, "turn this node into a seeder")
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()
	portStr := ":" + strconv.Itoa(*port)
	if *seed {
		seeder.Start("tcp", portStr)
	} else {
		n := node.NewNode()
		n.Start("tcp", portStr)
	}
}
