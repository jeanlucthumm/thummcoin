package main

import (
	"github.com/jeanlucthumm/thummcoin/cli"
	"github.com/jeanlucthumm/thummcoin/p2p"
	"os"
)

func main() {
	n := p2p.NewNode(false)
	cli.Interpret(os.Stdin, n)
}
