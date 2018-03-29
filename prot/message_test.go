package prot

import (
	"testing"
	"bytes"
)

func TestMakePingMessageAndDecode(t *testing.T) {
	m, err := MakePingMessage("me", "you")
	if err != nil {
		t.Error(err)
	}
	if m.ID != PING {
		t.Fail()
	}

	ptemp, err := DecodeMessage(m)
	if err != nil {
		t.Error(err)
	}
	p, ok := ptemp.(*Ping)
	if !ok {
		t.Fail()
	}

	if p.From != "me" || p.To != "you" {
		t.Fail()
	}
}

func TestMakePeerListMessageAndDecode(t *testing.T) {
	ips := []string{
		"192.168.1.1",
		"192.168.1.2",
		"192.168.1.3",
	}

	m, err := MakePeerListMessage(ips)
	if err != nil {
		t.Error(err)
	}
	if m.ID != PLIST {
		t.Fail()
	}

	ptemp, err := DecodeMessage(m)
	if err != nil {
		t.Error(err)
	}
	p, ok := ptemp.(*PeerList)
	if !ok {
		t.Fail()
	}

	for i, p := range p.Peers {
		if p.Address != ips[i] {
			t.Fail()
		}
	}
}

func TestSendReceive(t *testing.T) {
	ips := []string{
		"192.168.1.1",
		"192.168.1.2",
		"192.168.1.3",
	}
	m, _ := MakePeerListMessage(ips)

	var buf bytes.Buffer
	err := Send(&buf, m)
	if err != nil {
		t.Error(err)
	}

	mf, err := Receive(&buf)
	if err != nil {
		t.Error(err)
	}

	if bytes.Compare(m.Data, mf.Data) != 0 || m.ID != mf.ID {
		t.Fail()
	}
}
