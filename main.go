package main

import (
	"flag"
	"github.com/jeanlucthumm/thummcoin/seeder"
)

func main() {
	seed := flag.Bool("seed", false, "turn this node into a seeder")
	flag.Parse()
	if *seed {
		seeder.Start("tcp", ":8090")
	}
}
