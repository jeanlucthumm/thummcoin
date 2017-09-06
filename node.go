package main

import (
	"net"
	"log"
	"fmt"
	"encoding/gob"
)

type Node struct {
	peers map[*peer]bool // TODO use sync map
	ln    net.Listener

	broadcast chan *message // note consider making public
	add       chan *peer
	del       chan *peer
	stop      chan bool
}

func NewNode() *Node {
	return &Node{
		peers: make(map[*peer]bool),

		broadcast: make(chan *message),
		add:       make(chan *peer),
		del:       make(chan *peer),
		stop:      make(chan bool),
	}
}

func (n *Node) Start(addr net.Addr) error {
	var err error
	n.ln, err = net.Listen(addr.Network(), addr.String())
	if err != nil {
		log.Fatalln("Could not start node:", err)
		return err
	}
	log.Println("Node listening on:", addr)

	go n.handleChannels() // make node responsive
	go n.handleConnections()
	return nil
}

func (n *Node) Stop() {
	n.stop <- true
}

func (n *Node) handleChannels() {
	for {
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
		case <-n.stop:
			log.Println("Stopping server")
			n.ln.Close()
			for p := range n.peers {
				p.socket.Close()
			}
			break
		}
	}
}

func (n *Node) handleConnections() {
	for {
		conn, err := n.ln.Accept()
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
