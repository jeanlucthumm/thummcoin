package node

import (
	_ "fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jeanlucthumm/thummcoin/prot"
	"github.com/jeanlucthumm/thummcoin/util"
	"github.com/pkg/errors"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

const (
	p2pPort        = 8080
	ioTimeout      = time.Second * 2
	readBufferSize = 4096 // size of read buffer for incoming connections in bytes
)

// Node handles incoming connections and associated data
type Node struct {
	ln       net.Listener // listens for incoming connections
	peerList *peerList
	seed     bool
	ip       net.IPAddr

	broadcastChan chan *message
}

type message struct {
	kind prot.Type
	data []byte
}

// NewNode initializes a new Node but does not start it
func NewNode(seed bool) *Node {
	n := &Node{
		seed:          seed,
		broadcastChan: make(chan *message),
	}
	n.peerList = newPeerList(n)
	return n
}

// Start starts Node on addr, enabling it to respond to other nodes
func (n *Node) Start(addr net.Addr) error {
	var err error

	// instantiate listener
	n.ln, err = net.Listen(addr.Network(), addr.String())
	if err != nil {
		return errors.Wrap(err, "listen failed on node startup")
	}

	// make server responsive
	n.peerList.start()
	go n.handleChannels()
	if !n.seed {
		go n.discover()
	}

	for {
		conn, err := n.ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go n.handleConnection(conn)
	}
}

// discover attempts to find nodes and connect to the network. Must be called after Start.
// It does not check for self-connection and automatically dials seed, so it should not be used
// when in seed mode.
func (n *Node) discover() {
	// resolve seed
	conn, err := net.Dial("tcp", "seed:"+strconv.Itoa(p2pPort))
	if err != nil {
		log.Println("Failed to resolve seed host name")
		return
	}

	// identify IP
	reqIp := &prot.Request{Type: prot.Request_IP_SELF}
	riBuf, err := proto.Marshal(reqIp)
	if err != nil {
		log.Printf("Failed to marshal ip request during discovery: %s\n", err)
		return
	}
	mi := &message{
		kind: prot.Type_REQ,
		data: riBuf,
	}
	err = n.sendMessage(mi, conn)
	if err != nil {
		log.Printf("Failed to send ip req to seed: %s\n", err)
		return
	}

	ipPl := &prot.PeerList{}
	err = n.recvMessage(conn, ipPl)
	if err != nil {
		log.Printf("Failed to receive ip from seed: %s\n", err)
		return
	}
	if len(ipPl.Peers) == 0 {
		log.Println("Invalid self ip response from seed: peer list is empty")
		return
	}
	ip, err := net.ResolveIPAddr("ip", ipPl.Peers[0].Address)
	if err != nil {
		log.Printf("Failed to resolve ip response addr: %s\n", err)
		return
	}
	n.ip = *ip

	log.Printf("Self-identified as %s\n", n.ip.String())

	// request peer list
	reqPl := &prot.Request{Type: prot.Request_PEER_LIST}
	rplBuf, err := proto.Marshal(reqPl)
	if err != nil {
		log.Printf("Failed to marshal peer list req during discovery: %s\n", err)
		return
	}
	mpl := &message{
		kind: prot.Type_REQ,
		data: rplBuf,
	}
	err = n.sendMessage(mpl, conn)
	if err != nil {
		log.Printf("Failed to send peer list request to seed: %s\n", err)
		return
	}

	pl := &prot.PeerList{}
	err = n.recvMessage(conn, pl)
	if err != nil {
		log.Printf("Failed to receive peer list from seed: %s\n", err)
		return
	}

	err = conn.Close()
	if err != nil {
		log.Printf("Failed to close connection to seed: %s\n", err)
		// We continue anyways because that's seed's problem
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
	for {
		num, err := conn.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Failed read from %s: %s\n", remoteAddr)
			return
		}

		// register this peer
		if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
			n.peerList.newPeer <- util.IPFromTCP(*addr)
		}

		// route message
		msg := &prot.Message{}
		if err := proto.Unmarshal(b[:num], msg); err != nil {
			log.Printf("Failed to unmarshal message: %s\n", err)
		}

		switch msg.Type {
		case prot.Type_REQ:
			if err := n.handleRequest(conn, msg.Data); err != nil {
				log.Printf("Failed to handle request from %s: %s\n", remoteAddr, err)
				return
			}
		case prot.Type_PEER_LIST:
			// Seeds ignore peer lists
			if !n.seed {
				log.Printf("Got peer list from %s\n", remoteAddr)

			}
		}
	}
}

func (n *Node) broadcast(msg *message) {
	addrList := n.peerList.getAddresses()

	for _, ad := range addrList {
		conn, err := net.Dial("tcp", util.IPDialString(ad, p2pPort))
		if err != nil {
			log.Printf("Failed to dial %s during broadcast\n", ad)
			continue
		}

		err = n.sendMessage(msg, conn)
		if err != nil {
			log.Printf("Failed to send message to %s during broadcast: %s\n", ad, err)
			continue
		}

		err = conn.Close()
		if err != nil {
			log.Printf("Failed to close connection to %s during broadcast: %s\n", ad, err)
		}
	}
}
