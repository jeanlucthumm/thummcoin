package prot

import (
	"testing"
	"fmt"
	"net"
)

func TestTransaction_MarshalBinary(t *testing.T) {
	trans := Transaction{Amt: 22.3}
	copy(trans.Dest[:], "abcd")

	buf, err := trans.MarshalBinary()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	var o Transaction
	err = o.UnmarshalBinary(buf)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if trans != o {
		t.Fail()
	}
}

func TestPeerList_MarshalBinary(t *testing.T) {
	addr1, _ := net.ResolveTCPAddr("tcp", ":8080")
	addr2, _ := net.ResolveTCPAddr("tcp", ":8081")

	p := PeerList{
		Num:  2,
		List: []net.Addr{addr1, addr2},
	}

	buf, err := p.MarshalBinary()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	var o PeerList
	err = o.UnmarshalBinary(buf)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if p.Num != o.Num {
		fmt.Println("count did not match")
		t.Fail()
	}

	fmt.Println(p.List)
	fmt.Println(o.List)
}

func TestPing_MarshalBinary(t *testing.T) {
	p := Ping{
		From: ":8080",
		To:   ":8081",
	}

	buf, err := p.MarshalBinary()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	var o Ping
	err = o.UnmarshalBinary(buf)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if o != p {
		fmt.Print("o:", o, "p:", p)
		t.Fail()
	}
}

func TestPing_ValidateResponse(t *testing.T) {
	p := NewPing(":8080", ":8081")
	o := NewPing(":8081", ":8080")
	bad := NewPing(":200", ":8080")

	if !p.ValidateResponse(o) {
		t.Fail()
	}

	if p.ValidateResponse(bad) {
		t.Fail()
	}
}
