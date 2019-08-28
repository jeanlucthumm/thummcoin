package node

import (
	"fmt"
	_ "fmt"
	"github.com/pkg/errors"
	"log"
	"net"
	"sync"
	"time"
)

const (
	p2pPort       = 8080
	readDeadline  = 10 // read deadline in seconds for sockets
	writeDeadline = 10 // write deadline in seconds for sockets
)

// Node handles incoming connections and associated data
type Node struct {
	ln       net.Listener   // listens for incoming connections
	ptable   map[*peer]bool // look up table for all known peers
	tableMux sync.Mutex     // locks access to ptable. do not use in conjunction with channels

	addPeer chan *peer // adds a peer to ptable
	delPeer chan *peer // remove a peer from ptable
}

// peer represents a contactable peer
type peer struct {
	addr net.TCPAddr
}

// NewNode initializes a new Node but does not start it
func NewNode() *Node {
	return &Node{
		ptable:  make(map[*peer]bool),
		addPeer: make(chan *peer, 10),
		delPeer: make(chan *peer, 10),
	}
}

// Start starts Node on addr, enabling it to respond to other nodes
func (n *Node) Start(addr net.Addr) error {
	log.Println("Starting node")
	var err error

	// instantiate listener
	n.ln, err = net.Listen(addr.Network(), addr.String())
	if err != nil {
		return errors.Wrap(err, "listen failed on node startup")
	}

	// make server responsive
	go n.handleChannels()
	go n.Discover()
	go n.pingLoop()

	for {
		conn, err := n.ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go n.handleConnection(conn)
	}

	return nil
}

func (n *Node) StartSeed(addr net.Addr) error {
	log.Println("Starting seed")
	var err error

	// instantiate listener
	n.ln, err = net.Listen(addr.Network(), addr.String())
	if err != nil {
		return errors.Wrap(err, "listen failed on node startup")
	}

	// make server responsive
	go n.handleChannels()

	for {
		conn, err := n.ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go n.handleConnection(conn)
	}

	return nil
}

// Discover attempts to find nodes and connect to the network. Must be called after Start.
// It does not check for self-connection and automatically dials seed, so it should not be used
// when in seed mode.
func (n *Node) Discover() {
	// Dial the seed
	seedAddress := fmt.Sprintf("seed:%v", p2pPort)
	conn, err := net.Dial("tcp", seedAddress)
	if err != nil {
		log.Printf("Failed to dial seed at %s\n", seedAddress)
		return
	}

	// Test hello message
	if _, err := conn.Write([]byte("Hello there seed!")); err != nil {
		log.Println(errors.Wrap(err, "failed to ping seed").Error())
		return
	}
}

func (n *Node) handleChannels() {
	for {
		select {
		case p := <-n.addPeer:
			n.tableMux.Lock()
			n.ptable[p] = true
			n.tableMux.Unlock()
		case p := <-n.delPeer:
			n.tableMux.Lock()
			delete(n.ptable, p)
			n.tableMux.Unlock()
		}
	}
}

func (n *Node) handleConnection(conn net.Conn) {
	// TODO consider setting a read deadline
	b := make([]byte, 4096)
	num, err := conn.Read(b)
	if err != nil {
		log.Printf("Failed read from %s: %s\n", conn.RemoteAddr().String(), err.Error())
		return
	}

	log.Printf("From %s: %s\n", conn.RemoteAddr().String(), b[:num])

	if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		addr.Port = p2pPort
		n.addPeer <- &peer{addr: *addr}
	}
}

func (n *Node) pingLoop() {
	for {
		log.Println("Pinging all peers")
		n.pingAll()
		time.Sleep(5 * time.Second)
	}
}

func (n *Node) pingAll() {
	n.tableMux.Lock()
	defer n.tableMux.Unlock()

	for p := range n.ptable {
		// attempt to dial
		conn, err := net.Dial(p.addr.Network(), p.addr.String())
		if err != nil {
			log.Printf("Failed ping to %s: %s\n", p.addr.String(), err.Error())
			delete(n.ptable, p)
			continue
		}

		// attempt to write message
		_, err = conn.Write([]byte("Hello there, peer!"))
		if err != nil {
			log.Printf("Failed ping write to %s: %s\n", p.addr.String(), err.Error())
			delete(n.ptable, p)
			continue
		}
	}
}
