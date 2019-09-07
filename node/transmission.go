package node

import (
	"github.com/golang/protobuf/proto"
	"github.com/jeanlucthumm/thummcoin/prot"
	"github.com/jeanlucthumm/thummcoin/util"
	"github.com/pkg/errors"
	"log"
	"net"
	"time"
)

// transmission.go contains helper and handlers for different types of sent and received messages

func (n *Node) handleRequest(conn net.Conn, req *prot.Request) error {
	// identify request type and construct response data
	var kind prot.Type
	var buf []byte
	var err error
	switch req.Type {
	case prot.Request_PEER_LIST:
		log.Printf("Got peer list request from %s\n", conn.RemoteAddr())
		buf, err = n.makePeerList()
		if err != nil {
			return errors.Wrap(err, "make peer list")
		}
		kind = prot.Type_PEER_LIST
	case prot.Request_IP_SELF:
		log.Printf("Got ip request from %s\n", conn.RemoteAddr())
		buf, err = n.makeIpResponse(conn)
		if err != nil {
			return errors.Wrap(err, "make ip response")
		}
		kind = prot.Type_PEER_LIST
	default:
		return errors.New("unknown request type")
	}

	err = n.sendMessage(&Message{kind: kind, data: buf}, conn)
	if err != nil {
		return errors.Wrap(err, "send message")
	}

	return nil
}

func (n *Node) processPeerList(pl *prot.PeerList) {
	for _, p := range pl.Peers {
		ip, err := net.ResolveIPAddr("ip", p.Address)
		if err != nil {
			log.Printf("Failed to resolve ip address %s from peer list: %s\n", p.Address, err)
		}
		if util.IPEqual(*ip, n.ip) {
			continue
		}

		n.peerList.newPeer <- *ip
	}
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

func (n *Node) makeIpResponse(conn net.Conn) ([]byte, error) {
	tcp, err := net.ResolveTCPAddr("tcp", conn.RemoteAddr().String())
	if err != nil {
		return nil, errors.Wrap(err, "resolve tcp addr")
	}
	ip := net.IPAddr{
		IP:   tcp.IP,
		Zone: tcp.Zone,
	}
	p := &prot.PeerList_Peer{Address: ip.String()}
	pl := &prot.PeerList{Peers: []*prot.PeerList_Peer{p}}
	buf, err := proto.Marshal(pl)
	if err != nil {
		return nil, errors.Wrap(err, "marshal ip response")
	}
	return buf, nil
}

func (n *Node) sendMessage(msg *Message, conn net.Conn) error {
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
