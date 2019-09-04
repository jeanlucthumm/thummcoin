package node

import (
	"github.com/golang/protobuf/proto"
	"github.com/jeanlucthumm/thummcoin/prot"
	"github.com/pkg/errors"
	"log"
	"net"
	"time"
)

// transmission.go contains helper and handlers for different types of sent and received messages

const (
	ioTimeout = time.Second * 2
)

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
		log.Printf("Got peer list request from %s\n", conn.RemoteAddr().String())
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

func (n *Node) sendMessage(msg *message, conn net.Conn) error {
	m := &prot.Message{
		Type: msg.kind,
		From: n.ln.Addr().String(),
		To:   conn.RemoteAddr().String(),
		Data: msg.data,
	}
	buf, err := proto.Marshal(m)
	if err != nil {
		return errors.Wrap(err, "marshal message")
	}

	err = conn.SetDeadline(time.Now().Add(ioTimeout))
	if err != nil {
		return errors.Wrap(err, "set io deadline")
	}
	_, err = conn.Write(buf)
	if err != nil {
		return errors.Wrap(err, "write message to "+conn.RemoteAddr().String())
	}
	return nil
}

// pb is not message itself but the decoded data of the message
func (n *Node) recvMessage(conn net.Conn, pb proto.Message) error {
	buf := make([]byte, readBufferSize)
	err := conn.SetDeadline(time.Now().Add(ioTimeout))
	if err != nil {
		return errors.Wrap(err, "set io deadline")
	}
	num, err := conn.Read(buf)
	if err != nil {
		return errors.Wrap(err, "read from conn")
	}

	m := &prot.Message{}
	err = proto.Unmarshal(buf[:num], m)
	if err != nil {
		return errors.Wrap(err, "unmarshal message")
	}

	// TODO verify that it was meant for us
	//  	- need better ways of comparing IPs

	err = proto.Unmarshal(m.Data, pb)
	if err != nil {
		return errors.Wrap(err, "unmarshal message data")
	}

	return nil
}
