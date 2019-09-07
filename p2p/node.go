package p2p

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
	seedTimeout    = 2
	readBufferSize = 4096 // size of read buffer for incoming connections in bytes
)

// Node handles incoming connections and associated data
type Node struct {
	ln       net.Listener // listens for incoming connections
	peerList *peerList
	seed     bool
	ip       net.IPAddr

	Broadcast chan *Message
	Done      chan bool
}

type Message struct {
	Kind prot.Type
	Data []byte
}

// NewNode initializes a new Node but does not start it
func NewNode(seed bool) *Node {
	n := &Node{
		seed:      seed,
		Broadcast: make(chan *Message),
		Done:      make(chan bool),
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

	go func() {
		for {
			conn, err := n.ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}

			go n.handleConnection(conn)
		}
	}()

	return nil
}

// discover attempts to find nodes and connect to the network. Must be called after Start.
// It does not check for self-connection and automatically dials seed, so it should not be used
// when in seed mode.
func (n *Node) discover() {
	// resolve seed
	var conn net.Conn
	for {
		var err error
		conn, err = net.Dial("tcp", "seed:"+strconv.Itoa(p2pPort))
		if err == nil {
			break
		}
		log.Println("Failed to resolve seed host name")
		time.Sleep(time.Second * seedTimeout)
	}

	// identify IP
	reqIp := &prot.Request{Type: prot.Request_IP_SELF}
	riBuf, err := proto.Marshal(reqIp)
	if err != nil {
		log.Printf("Failed to marshal ip request during discovery: %s\n", err)
		return
	}
	mi := &Message{
		Kind: prot.Type_REQ,
		Data: riBuf,
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
	mpl := &Message{
		Kind: prot.Type_REQ,
		Data: rplBuf,
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

	n.processPeerList(pl)
}

func (n *Node) handleChannels() {
	for {
		select {
		case msg := <-n.Broadcast:
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
			continue
		}

		switch msg.Type {
		case prot.Type_REQ:
			req := &prot.Request{}
			err = proto.Unmarshal(msg.Data, req)
			if err != nil {
				log.Printf("Failed to unmarshal request from %s: %s\n", remoteAddr, err)
				continue
			}
			if err := n.handleRequest(conn, req); err != nil {
				log.Printf("Failed to handle request from %s: %s\n", remoteAddr, err)
				continue
			}
		case prot.Type_PEER_LIST:
			// Seeds ignore peer lists
			if n.seed {
				continue
			}
			log.Printf("Got peer list from %s\n", remoteAddr)
			pl := &prot.PeerList{}
			err = proto.Unmarshal(msg.Data, pl)
			if err != nil {
				log.Printf("Failed to unmarshal peer list from %s: %s\n", remoteAddr, err)
				continue
			}
			go n.processPeerList(pl)
		case prot.Type_TEXT:
			log.Printf("Got text from %s: %s\n", remoteAddr, string(msg.Data))
		}
	}
}

func (n *Node) broadcast(msg *Message) {
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
