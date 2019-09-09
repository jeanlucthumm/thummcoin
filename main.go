package main

import (
	"flag"
	"github.com/jeanlucthumm/thummcoin/cli"
	"github.com/jeanlucthumm/thummcoin/p2p"
	"github.com/sirupsen/logrus"
	"log"
	"net"
	"os"
)

func main() {
	logrus.SetLevel(logrus.TraceLevel)

	var seedMode = flag.Bool("seed", false, "enable seeding mode")
	flag.Parse()

	n := p2p.NewNode(*seedMode)
	addr, _ := net.ResolveTCPAddr("tcp", ":8080")

	log.Println()

	if err := n.Start(addr); err != nil {
		log.Printf("Failed to start node: %s\n", err)
		return
	}

	cli.Interpret(os.Stdin, n)
}
