package node

import (
	"fmt"
	_ "fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jeanlucthumm/thummcoin/prot"
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

	readBufferSize = 4096 // size of read buffer for incoming connections in bytes
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
	// dial the seed
	seedAddress := fmt.Sprintf("seed:%v", p2pPort)
	conn, err := net.Dial("tcp", seedAddress)
	if err != nil {
		log.Printf("Failed to dial seed at %s\n", seedAddress)
		return
	}

	// request peer list
	req := &prot.Request{Type: prot.Request_PEER_LIST}
	buf, err := proto.Marshal(req);
	if err != nil {
		log.Printf("Failed to marshal request during discovery: %s\n", err.Error())
	}
	m := &prot.Message{
		Type: prot.Type_REQ,
		From: n.ln.Addr().String(),
		To:   "seed",
		Data: buf,
	}
	mBuf, err := proto.Marshal(m);
	if err != nil {
		log.Printf("Failed to marshal message during discovery: %s\n", err.Error())
	}
	if _, err = conn.Write(mBuf); err != nil {
		log.Printf("Failed to write request to seed: %s\n", err.Error())
	}

	// Wait for response
	num, err := conn.Read(buf)
	if err != nil {
		log.Printf("Failed to read back from seed during discovery: %s\n", err)
	}
	log.Printf("Read %d bytes from seed\n", num)
}

// TODO the port num is changed. need to keep session alive as we wait for response in client.
// 		Other option is to only do stateless communication, but that defeats the point of TCP

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
	remoteAddr := conn.RemoteAddr().String()
	b := make([]byte, readBufferSize) // FIXME messages can be much larger than that
	num, err := conn.Read(b)
	if err != nil {
		log.Printf("Failed read from %s: %s\n", remoteAddr)
		return
	}

	log.Printf("Read %d bytes from %s\n", num, remoteAddr)

	// register this peer
	if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		addrCopy := *addr
		addrCopy.Port = p2pPort
		n.addPeer <- &peer{addr: addrCopy}
	}

	// route message
	msg := &prot.Message{}
	if err := proto.Unmarshal(b[:num], msg); err != nil {
		log.Printf("Failed to unmarshal message: %s\n", err)
	}

	switch msg.Type {
	case prot.Type_REQ:
		log.Println("Got request")
		if err := n.handleRequest(conn, msg.Data); err != nil {
			log.Printf("Failed to handle request from %s: %s\n", remoteAddr, err)
		}
	case prot.Type_PEER_LIST:
		log.Println("Got peer list")
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
	// FIXME this mux locks the table for too long. Get a copy of all IPs instead

	for p := range n.ptable {
		// attempt to dial
		conn, err := net.Dial(p.addr.Network(), p.addr.String())
		if err != nil {
			log.Printf("Failed ping to %s: %s\n", p.addr.String(), err)
			delete(n.ptable, p)
			continue
		}

		// attempt to write message
		_, err = conn.Write([]byte("Hello there, peer!"))
		if err != nil {
			log.Printf("Failed ping write to %s: %s\n", p.addr.String(), err)
			delete(n.ptable, p)
			continue
		}
	}
}
