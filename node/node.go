package node

import (
	"fmt"
	_ "fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jeanlucthumm/thummcoin/prot"
	"github.com/jeanlucthumm/thummcoin/util"
	"github.com/pkg/errors"
	"log"
	"net"
)

const (
	p2pPort       = 8080
	readDeadline  = 10 // read deadline in seconds for sockets
	writeDeadline = 10 // write deadline in seconds for sockets

	readBufferSize = 4096 // size of read buffer for incoming connections in bytes
)

// Node handles incoming connections and associated data
type Node struct {
	ln       net.Listener // listens for incoming connections
	peerList *peerList

	broadcastChan chan []byte
}

// peer represents a contactable peer
type Peer struct {
	addr net.IP
}

// NewNode initializes a new Node but does not start it
func NewNode() *Node {
	return &Node{
		peerList:      newPeerList(),
		broadcastChan: make(chan []byte),
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
	n.peerList.start()
	go n.handleChannels()
	go n.Discover()

	for {
		conn, err := n.ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go n.handleConnection(conn)
	}
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
	rBuf, err := proto.Marshal(req)
	if err != nil {
		log.Printf("Failed to marshal request during discovery: %s\n", err.Error())
	}
	m := &prot.Message{
		Type: prot.Type_REQ,
		From: n.ln.Addr().String(),
		To:   "seed",
		Data: rBuf,
	}
	mBuf, err := proto.Marshal(m)
	if err != nil {
		log.Printf("Failed to marshal message during discovery: %s\n", err.Error())
	}
	if _, err = conn.Write(mBuf); err != nil {
		log.Printf("Failed to write request to seed: %s\n", err.Error())
	}

	// wait for response
	buf := make([]byte, readBufferSize)
	num, err := conn.Read(buf)
	if err != nil {
		log.Printf("Failed to read back from seed during discovery: %s\n", err)
	}

	// decode and all known peers
	pl := &prot.PeerList{}
	err = proto.Unmarshal(buf[:num], pl)
	if err != nil {
		log.Printf("Failed to unmarshal peer list from seed: %s\n", err)
	}

	for _, peer := range pl.Peers {
		log.Printf("Adding peer with address %s\n", peer.Address)
	}
}

func (n *Node) handleChannels() {
	for {
		select {
		case msg := <-n.broadcastChan:
			go n.broadcast(msg)
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
		n.peerList.newPeer <- addr.IP
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
		log.Println("Got peer list") // TODO
	}
}

func (n *Node) broadcast(msg []byte) {
	addrList := n.peerList.getAddresses()

	for _, ad := range addrList {
		conn, err := net.Dial("tcp", util.IPString(ad, p2pPort))
		if err != nil {
			// TODO do some sort of retry procedure then drop peer
			continue
		}

		_, err = conn.Write(msg)
		if err != nil {
			log.Printf("Failed to write msg to %s during broadcast: %s\n",
				ad.String(), err)
			continue
		}
	}
}
