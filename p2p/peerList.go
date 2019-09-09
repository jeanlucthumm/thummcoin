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

	newPeers chan []*net.IPAddr
	stop     chan bool
}

type peer struct {
	addr net.IPAddr
}

func newPeerList(node *Node) *peerList {
	return &peerList{
		node:     node,
		newPeers: make(chan []*net.IPAddr, peerBufferCap),
		stop:     make(chan bool),
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
		case addr := <-pl.newPeers:
			go pl.handleNewPeers(addr)
		case <-pl.stop:
			return
		}
	}
}

func (pl *peerList) handleNewPeers(addrs []*net.IPAddr) {
	// Since seed nodes are for minimal bootstrap, no need for broadcast
	newAddrs := pl.addNewAddrs(addrs)
	if len(newAddrs) != 0 && !pl.node.seed {
		// broadcast new peers to everyone
		var peers []*prot.PeerList_Peer
		for _, a := range newAddrs {
			peers = append(peers, &prot.PeerList_Peer{Address: a.String()})
		}
		plm := &prot.PeerList{Peers: peers}
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

func (pl *peerList) addNewAddrs(addrs []*net.IPAddr) []*net.IPAddr {
	pl.mux.Lock()
	defer pl.mux.Unlock()

	var newAddrs []*net.IPAddr

	for _, addr := range addrs {
		seen := false
		for _, p := range pl.list {
			if util.IPEqual(addr, &p.addr) {
				seen = true
			}
		}
		if !seen {
			newAddrs = append(newAddrs, addr)
			pl.list = append(pl.list, newPeer(addr))
		}
	}

	return newAddrs
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
