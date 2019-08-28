package prot

import (
	"testing"
	"bytes"
	"strings"
)

func TestMakePingMessageAndDecode(t *testing.T) {
	m, err := MakePingMessage(&Ping{From: "me", To: "you"})
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
	m, err := MakePeerListMessage(ips)
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = Send(&buf, m)
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

	// bad cases
	m, err = Receive(strings.NewReader("Hello World"))
	if err == nil {
		t.Fail()
	}
}

func TestMakeMessage(t *testing.T) {
	// PING
	p := Ping{From: "me", To: "you"}
	m, err := MakeMessage(p)
	if err != nil {
		t.Error(err)
	}

	r, err := DecodeMessage(m)
	if err != nil {
		t.Error(err)
	}
	if r, ok := r.(*Ping); ok {
		if r.From != p.From || r.To != p.To {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}

func TestSample(t *testing.T) {
	//p := Ping{From: "me", To: "you"}
	//var i interface{}
	//i = p
	//if _, ok := i.(proto.Message); !ok {
	//	t.Fail()
	//}
}
