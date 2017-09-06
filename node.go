package main

import (
	"net"
	"sync"
	"log"
)

type Node struct {
	peers map[net.Addr]bool
	mux   sync.Mutex
}

func (n *Node) Start(addr net.Addr) error {
	n.mux.Lock()
	defer n.mux.Unlock()

	ln, err := net.Listen(addr.Network(), addr.String())
	if err != nil {
		log.Fatalln("Could not start node:", err)
	}
	log.Println("Node listening on:", addr)

	for {
		conn, _ := ln.Accept()
		peer := &Peer{socket: conn}
		go peer.handle()
	}
}

type Peer struct {
	socket net.Conn
}

func (p *Peer) handle() {
}
