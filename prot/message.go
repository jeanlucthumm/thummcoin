package prot

import (
	"io"
	"encoding/gob"
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
