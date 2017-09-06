package main

import (
	"net"
	"sync"
	"log"
	"fmt"
	"encoding/gob"
)

type Node struct {
	peers map[*peer]bool
	mux   sync.Mutex

	broadcast chan *message // note consider making public
	add       chan *peer
	del       chan *peer
}

func NewNode() *Node {
	return &Node{
		peers: make(map[*peer]bool),

		broadcast: make(chan *message),
		add:       make(chan *peer),
		del:       make(chan *peer),
	}
}

func (n *Node) Start(addr net.Addr) error {
	n.mux.Lock()
	defer n.mux.Unlock()

	ln, err := net.Listen(addr.Network(), addr.String())
	if err != nil {
		log.Fatalln("Could not start node:", err)
		return err
	}
	log.Println("Node listening on:", addr)

	go n.handleChannels() // make node responsive
	go n.handleConnections(ln)
	return nil
}

func (n *Node) handleChannels() {
	select {
	case msg := <-n.broadcast:
		for p := range n.peers {
			// encode the message and send
			enc := p.encoder()
			err := enc.Encode(msg) // Q a way to store encoded form instead of doing it every time?
			if err != nil {
			}
		}
	case p := <-n.add:
		n.peers[p] = true
	}
}

func (n *Node) handleConnections(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		fmt.Println("New connection:", conn.RemoteAddr()) // DEBUG
		peer := &peer{socket: conn}
		go peer.Handle()
	}
}

type peer struct {
	socket net.Conn
	node   *Node
}

func (p *peer) Handle() {
	dec := p.decoder()
	for {
		msg := new(message)
		err := dec.Decode(msg)
		if err != nil {
			fmt.Println("Could not read from peer:", p.socket.RemoteAddr(), err) // DEBUG
		}
	}
}

func (p *peer) Write(msg message) error {
	enc := p.encoder()
	err := enc.Encode(msg)
	if err != nil {
		fmt.Println("Could not write to peer:", p.socket.RemoteAddr(), err) // DEBUG
	}
	return err
}

func (p *peer) encoder() *gob.Encoder {
	// TODO store this as a field
	return gob.NewEncoder(p.socket)
}

func (p *peer) decoder() *gob.Decoder {
	// TODO store this as a field
	return gob.NewDecoder(p.socket)
}

func (p *peer) end() {
	p.socket.Close()
}
