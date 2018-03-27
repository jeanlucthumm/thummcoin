package node

import "github.com/jeanlucthumm/thummcoin/prot"

func makePeerListMessage(ptable map[*peer]bool) (*prot.Message, error){
	m := &prot.Message{}
	pl := &prot.PeerList{}
	// populate prot peer list
	for p := range ptable {
		pl.Peers = append(pl.Peers, &prot.PeerList_Peer{
			Address: p.addr.IP.String(),
		})
	}
}
