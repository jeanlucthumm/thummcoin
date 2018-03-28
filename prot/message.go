package prot

import (
	"io"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
)

const (
	PING  = iota
	PLIST
)

type Message struct {
	ID   int
	Data []byte
}

// Send encodes the message and writes it to w
func Send(w io.Writer, message *Message) error {
	e := gob.NewEncoder(w)
	return e.Encode(*message)
}

// Receive reads from r and extracts a message
func Receive(r io.Reader) (*Message, error) {
	var m *Message
	d := gob.NewDecoder(r)
	err := d.Decode(m)
	if err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func DecodeMessage(m *Message) (interface{}, error) {
	switch m.ID {
	case PING:

	}
}

// MakePeerListMessage constructs a peer list message from the given ips
func MakePeerListMessage(ips []string) (*Message, error) {
	m := &Message{}erm
	pl := &PeerList{}

	// populate peer list
	for _, ip := range ips {
		pl.Peers = append(pl.Peers, &PeerList_Peer{
			Address: ip,
		})
	}

	// construct message
	var err error
	m.ID = PLIST
	m.Data, err = proto.Marshal(pl)
	if err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

// MakePingMessage constructs a ping message from the given parameters
func MakePingMessage(from, to string) (*Message, error) {
	m := &Message{}
	p := &Ping{}
	p.From = from
	p.To = to

	var err error
	m.ID = PING
	m.Data, err = proto.Marshal(p)
	if err != nil {
		return nil, err
	} else {
		return m, nil
	}
}
