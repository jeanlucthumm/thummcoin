package node

import (
	"net"
	"log"
	"encoding/gob"
	"github.com/jeanlucthumm/thummcoin/seeder"
	"github.com/jeanlucthumm/thummcoin/prot"
	"time"
	"fmt"
)

/// Node

const (
	// time out for reading from other nodes
	rtimeout = 20 * time.Second
)

type message struct {
	id   byte
	data []byte
}

type Node struct {
	peers map[*peer]bool // TODO use sync map
	ln    net.Listener

	Broadcast chan *message
	add       chan *peer
	del       chan *peer
	stop      chan bool
}

func NewNode() *Node {
	return &Node{
		peers: make(map[*peer]bool),

		Broadcast: make(chan *message),
		add:       make(chan *peer),
		del:       make(chan *peer),
		stop:      make(chan bool),
	}
}

func (n *Node) Start(network string, loc string) {
	addr, err := net.ResolveTCPAddr(network, loc)
	if err != nil {
		log.Fatalln("Could not start node:", err)
	}

	n.ln, err = net.Listen(addr.Network(), addr.String())
	if err != nil {
		log.Fatalln("Could not start node:", err)
	}

	go n.handleChannels() // make node responsive
	go n.discoverPeers()
	n.handleConnections()
}

func (n *Node) Stop() {
	n.stop <- true
}

func (n *Node) handleChannels() {
	for {
		select {
		case msg := <-n.Broadcast:
			for p := range n.peers {
				// encode the message and send
				enc := p.encoder()
				err := enc.Encode(msg) // Q a way to store encoded form instead of doing it every time?
				if err != nil {
				}
			}
		case p := <-n.add:
			log.Println("New peer:", p.socket.RemoteAddr())
			n.peers[p] = true
		case p := <-n.del:
			log.Println("Deleting peer:", p.socket.RemoteAddr())
			delete(n.peers, p)
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
	log.Println("Node listening on:", n.ln.Addr())
	for {
		// FIXME what if a peer reconnects?
		conn, err := n.ln.Accept()
		if err != nil {
			continue
		}
		log.Println("New connection:", conn.RemoteAddr()) // DEBUG
		peer := &peer{socket: conn, node: n}
		go peer.Handle()
	}
}

func (n *Node) discoverPeers() {
	for _, seedIP := range seeder.SeederIPs {
		seed, err := net.Dial("tcp", seedIP)
		if err != nil {
			log.Println("Dead seeder:", seedIP)
			continue
		}

		pbuf := make([]byte, 4096) // FIXME we can't predict this size
		seed.SetReadDeadline(time.Now().Add(rtimeout))
		_, err = seed.Read(pbuf)
		if err != nil {
			log.Println("Bad seeder:", seedIP)
			continue
		}

		var plist prot.PeerList
		err = plist.UnmarshalBinary(pbuf)
		if err != nil {
			continue
		}

		for _, addr := range plist.List {
			go n.buddy(addr)
		}
		seed.Close()
	}
}

func (n *Node) buddy(addr net.Addr) {
	// FIXME make sure that you do not buddy with yourself

	conn, err := net.Dial(addr.Network(), addr.String())
	if err != nil {
		return
	}

	ping := prot.NewPing(n.ln.Addr().String(), conn.RemoteAddr().String())
	pingBuf, _ := ping.MarshalBinary()
	conn.Write(pingBuf)

	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(rtimeout))
	_, err = conn.Read(buf)
	if err != nil {
		log.Println("Bad buddy:", addr)
		conn.Close()
		return
	}

	var o prot.Ping
	err = o.UnmarshalBinary(buf)
	if err != nil || !ping.ValidateResponse(o) {
		log.Println("Bad buddy:", addr)
		conn.Close()
		return
	}

	p := &peer{socket: conn, node: n}
	n.add <- p
}

/// peer

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
			log.Printf("Could not read from peer %v: %v\n", p.socket.RemoteAddr(), err)
			p.node.del <- p
			break
		}
	}
}

func (p *peer) Write(msg message) error {
	enc := p.encoder()
	err := enc.Encode(msg)
	if err != nil {
		log.Println("Could not write to peer:", p.socket.RemoteAddr(), err) // DEBUG
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
