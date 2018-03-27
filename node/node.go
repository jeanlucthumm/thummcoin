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
	tableMux sync.Mutex     // locks access to ptable

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
		addPeer: make(chan *peer, 10),
		delPeer: make(chan *peer, 10),
	}
}

// Start starts Node on addr, enabling it to respond to other nodes
func (n *Node) Start(addr net.Addr) error {
	var err error

	// instantiate listener
	n.ln, err = net.Listen(addr.Network(), addr.String())
	if err != nil {
		return err
	}

	// make server responsive
	go n.handleChannels()
	go n.listen()
	go n.pingLoop()
	return nil
}

func (n *Node) listen() {
	for {
		conn, err := n.ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go n.handleConnection(conn)
	}
}

func (n *Node) handleChannels() {
	for {
		select {
		case p := <-n.addPeer:
			n.ptable[p] = true
		case p := <-n.delPeer:
			delete(n.ptable, p)
		}
	}
}

func (n *Node) handleConnection(conn net.Conn) {
	// TODO consider setting a read deadline
	var b []byte
	num, err := conn.Read(b)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("From %s: %s\n", conn.RemoteAddr().String(), b[:num])
}

func (n *Node) pingLoop() {
	for {
		log.Println("Pinging all peers")
		n.pingAll()
		time.Sleep(5 * time.Second)
	}
}

func (n *Node) pingAll() {
	// DEBUG
	n.tableMux.Lock() // TODO this is a huge lockout, fix it
	defer n.tableMux.Unlock()

	var delList []*peer
	for p, _ := range n.ptable {
		// attempt to dial
		conn, err := net.Dial(p.addr.Network(), p.addr.String())
		if err != nil {
			log.Println(err)
			delList = append(delList, p)
			continue
		}

		// attempt to write message
		_, err = conn.Write([]byte("Hello there, peer!"))
		if err != nil {
			log.Println(err)
			delList = append(delList, p)
			continue
		}
	}

	// remove bad peers
	for _, d := range delList {
		n.delPeer <- d
	}
}
