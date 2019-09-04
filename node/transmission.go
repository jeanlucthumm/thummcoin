package node

import (
	"github.com/golang/protobuf/proto"
	"github.com/jeanlucthumm/thummcoin/prot"
	"github.com/pkg/errors"
	"log"
	"net"
)

// transmission.go contains helper and handlers for different types of sent and received messages

func (n *Node) handleRequest(conn net.Conn, data []byte) error {
	// recover request
	req := &prot.Request{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		return errors.Wrap(err, "unmarshal request")
	}

	// identify request type and construct response data
	var dType prot.Type
	var buf []byte
	switch req.Type {
	case prot.Request_PEER_LIST:
		buf, err = n.makePeerList()
		if err != nil {
			return errors.Wrap(err, "make peer list")
		}
		dType = prot.Type_PEER_LIST
	default:
		return errors.New("unknown request type")
	}

	// wrap data in message and send
	m := &prot.Message{
		Type: dType,
		From: n.ln.Addr().String(),
		To:   conn.RemoteAddr().String(),
		Data: buf,
	}

	mBuf, err := proto.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "marshal message")
	}
	_, err = conn.Write(mBuf)
	if err != nil {
		return errors.Wrap(err, "transmit message")
	}
	return nil
}

func (n *Node) makePeerList() ([]byte, error) {
	addrList := n.peerList.getAddresses()

	pl := &prot.PeerList{}
	for _, ad := range addrList {
		p := &prot.PeerList_Peer{Address: ad.String()}
		pl.Peers = append(pl.Peers, p)
	}

	buf, err := proto.Marshal(pl)
	if err != nil {
		return nil, errors.Wrap(err, "marshal peer list")
	}
	return buf, nil
}
