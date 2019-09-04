package node

import (
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

func newPeerList() *peerList {
	return &peerList{
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
		// TODO broadcast new peer to everyone
	}
}

func (pl *peerList) addAddrIfNew(addr net.IP) bool {
	pl.mux.Lock()
	defer pl.mux.Unlock()

	found := false
	for _, p := range pl.list {
		if addr.Equal(p.addr) {
			found = true
		}
	}

	if !found {
		pl.list = append(pl.list, newPeer(addr))
	}
	return found
}

func (pl *peerList) getAddresses() []net.IP {
	pl.mux.Lock()
	defer pl.mux.Unlock()
	list := make([]net.IP, len(pl.list))
	for _, p := range pl.list {
		list = append(list, p.addr)
	}
	return list
}
