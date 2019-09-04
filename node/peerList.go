package node

import (
	"github.com/golang/protobuf/proto"
	"github.com/jeanlucthumm/thummcoin/prot"
	"log"
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

	newPeer chan net.IP
	stop    chan bool
}

type peer struct {
	addr net.IP
}

func newPeerList(node *Node) *peerList {
	return &peerList{
		node:    node,
		newPeer: make(chan net.IP, peerBufferCap),
		stop:    make(chan bool),
	}
}

func newPeer(addr net.IP) *peer {
	return &peer{
		addr: addr,
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

func (pl *peerList) handleNewPeer(addr net.IP) {
	if pl.addAddrIfNew(addr) {
		// broadcast new peer to everyone
		p := &prot.PeerList_Peer{Address: addr.String()}
		plm := &prot.PeerList{Peers: []*prot.PeerList_Peer{p}}
		buf, err := proto.Marshal(plm)
		if err != nil {
			log.Printf("Failed to marshal peer list for new peer: %s\n", err)
			return
		}
		pl.node.broadcastChan <- &message{
			kind: prot.Type_PEER_LIST,
			data: buf,
		}
	}
}

func (pl *peerList) addAddrIfNew(addr net.IP) bool {
	pl.mux.Lock()
	defer pl.mux.Unlock()

	newAd := true
	for _, p := range pl.list {
		if addr.Equal(p.addr) {
			newAd = false
		}
	}

	if newAd {
		pl.list = append(pl.list, newPeer(addr))
	}
	return newAd
}

func (pl *peerList) getAddresses() []net.IP {
	pl.mux.Lock()
	defer pl.mux.Unlock()
	list := make([]net.IP, len(pl.list))
	for i, p := range pl.list {
		list[i] = p.addr
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
