package p2p

import (
	"github.com/golang/protobuf/proto"
	"github.com/jeanlucthumm/thummcoin/prot"
	"github.com/jeanlucthumm/thummcoin/util"
	"net"
	"sync"
)

const (
	peerBufferCap = 50
)

type peerList struct {
	mux sync.Mutex
	// FIXME linear search is not efficient
	list []*peer
	node *Node

	newPeer chan *net.IPAddr
	stop    chan bool
}

type peer struct {
	addr net.IPAddr
}

func newPeerList(node *Node) *peerList {
	return &peerList{
		node:    node,
		newPeer: make(chan *net.IPAddr, peerBufferCap),
		stop:    make(chan bool),
	}
}

func newPeer(addr *net.IPAddr) *peer {
	return &peer{
		addr: *addr,
	}
}

func (pl *peerList) start() {
	go pl.handleChannels()
}

func (pl *peerList) handleChannels() {
	for {
		select {
		case addr := <-pl.newPeer:
			go pl.handleNewPeer(addr)
		case <-pl.stop:
			return
		}
	}
}

func (pl *peerList) handleNewPeer(addr *net.IPAddr) {
	// Since seed nodes are for minimal bootstrap, no need for broadcast
	if pl.addAddrIfNew(addr) && !pl.node.seed {
		// broadcast new peer to everyone
		p := &prot.PeerList_Peer{Address: addr.String()}
		plm := &prot.PeerList{Peers: []*prot.PeerList_Peer{p}}
		buf, err := proto.Marshal(plm)
		if err != nil {
			glog.Errorf("Failed to marshal peer list for new peer: %s", err)
			return
		}
		pl.node.Broadcast <- &Message{
			Kind: prot.Type_PEER_LIST,
			Data: buf,
		}
	}
}

func (pl *peerList) addAddrIfNew(addr *net.IPAddr) bool {
	pl.mux.Lock()
	defer pl.mux.Unlock()

	newAd := true
	for _, p := range pl.list {
		if util.IPEqual(addr, &p.addr) {
			newAd = false
		}
	}

	if newAd {
		pl.list = append(pl.list, newPeer(addr))
	}
	return newAd
}

func (pl *peerList) getAddresses() []*net.IPAddr {
	pl.mux.Lock()
	defer pl.mux.Unlock()
	list := make([]*net.IPAddr, len(pl.list))
	for i, p := range pl.list {
		list[i] = &p.addr
	}
	return list
}

func (pl *peerList) String() string {
	str := "peerList{"
	for i, p := range pl.list {
		if i != 0 {
			str += "," + p.addr.String()
		} else {
			str += p.addr.String()
		}
	}
	str += "}"
	return str
}
