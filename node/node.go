package node

import (
	"net"
	"log"
	"sync"
	"time"
)

// Node handles incoming connections and associated data
type Node struct {
	ln       net.Listener   // listens for incoming connections
	ptable   map[*peer]bool // look up table for all known peers
	tableMux sync.Mutex     // locks access to ptable. do not use in conjunction with channels

	addPeer chan *peer
	delPeer chan *peer
}

// peer represents a contactable peer
type peer struct {
	addr net.Addr
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
		return err
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

	n.ln, err = net.Listen(addr.Network(), addr.String())
	if err != nil {
		return err
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
	// dial seed node
	conn, err := net.Dial("tcp", "seed:8080")
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = conn.Write([]byte("Hello seed"))
	if err != nil {
		log.Fatal(err)
		return
	}

	n.addPeer <- &peer{addr: conn.RemoteAddr()}
}

func (n *Node) listen() {
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
		log.Println(err)
		return
	}

	log.Printf("From %s: %s\n", conn.RemoteAddr().String(), b[:num])

	n.addPeer <- &peer{addr: conn.RemoteAddr()}
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
			log.Println(err)
			delete(n.ptable, p)
			continue
		}

		// attempt to write message
		_, err = conn.Write([]byte("Hello there, peer!"))
		if err != nil {
			log.Println(err)
			delete(n.ptable, p)
			continue
		}
	}
}
