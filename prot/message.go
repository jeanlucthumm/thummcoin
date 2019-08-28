package prot

import (
	"io"
	"encoding/gob"
	"github.com/golang/protobuf/proto"
	"errors"
	"bytes"
)

const (
	PING  = iota
	PLIST
	IPREQ
)

type Message struct {
	ID   int
	Data []byte
}

func (m *Message) Equal(o *Message) bool {
	return o.ID == m.ID && bytes.Equal(o.Data, m.Data)
}

// Send encodes the message and writes it to w
func Send(w io.Writer, message *Message) error {
	e := gob.NewEncoder(w)
	return e.Encode(*message)
}

// Receive reads from r and extracts a message
func Receive(r io.Reader) (*Message, error) {
	m := &Message{}
	d := gob.NewDecoder(r)
	err := d.Decode(m)
	if err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

// DecodeMessage converts m to a concrete message. The result is returned as an empty interface
// and must be checked for type.
func DecodeMessage(m *Message) (interface{}, error) {
	switch m.ID {
	case PING:
		ping := &Ping{}
		err := proto.Unmarshal(m.Data, ping)
		return ping, err
	case PLIST:
		plist := &PeerList{}
		err := proto.Unmarshal(m.Data, plist)
		return plist, err
	default:
		return nil, errors.New("unknown message type")
	}
}

// MakePeerListMessage constructs a peer list message from the given ips
func MakePeerListMessage(ips []string) (*Message, error) {
	m := &Message{}
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

func MakePingMessage(p *Ping) (*Message, error){
	m := &Message{}

	var err error
	m.ID = PING
	m.Data, err = proto.Marshal(p)
	if err != nil {
		return nil, err
	} else {
		return m, nil
	}
}

func MakeMessage(e interface{}) (*Message, error) {
	m := &Message{}

	switch e.(type) {
	case Ping:
		m.ID = PING
	case PeerList:
		m.ID = PLIST
	case IPReq:
		m.ID = IPREQ
	default:
		return nil, errors.New("could not make message of unknown type")
	}

	if e, ok := e.(proto.Message); ok {
		var err error
		m.Data, err = proto.Marshal(e)
		if err != nil {
			return nil, err
		}
	}
	return nil, errors.New("could not type cast to proto.Message to marshal")
}
