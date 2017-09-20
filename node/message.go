package node

import (
	"bytes"
	"net"
	"encoding/binary"
	"fmt"
	"bufio"
	"strings"
)

// TODO move everything but essage to a protocol package

const (
	TRANS = 0x01
)

type message struct {
	id   byte
	data []byte
}

type Transaction struct {
	Dest [4]byte
	Amt  float64
}

func (t *Transaction) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, t)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *Transaction) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)
	return binary.Read(buf, binary.BigEndian, t)
}

type PeerList struct {
	Num  uint32 // TODO make these private since they are linked
	List []net.Addr
}

func (p *PeerList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, p.Num)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, addr := range p.List {
		buf.Write([]byte(addr.String() + "\n"))
	}
	return buf.Bytes(), nil
}

func (p *PeerList) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.BigEndian, &p.Num)
	if err != nil {
		fmt.Println(err)
		return err
	}

	buf := bufio.NewReader(r)
	for i := uint32(0); i < p.Num; i++ {
		str, err := buf.ReadString('\n')
		if err != nil {
			return err
		}
		addr, err := net.ResolveTCPAddr("tcp", strings.Trim(str, "\n"))
		if err != nil {
			return err
		}
		p.List = append(p.List, addr)
	}
	return nil
}

// Ping is used to communicate the presence of nodes
type Ping struct {
	From string
	To   string
}

func (p *Ping) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	_, err := buf.Write([]byte(p.From + "\n"))
	if err != nil {
		return nil, err
	}
	_, err = buf.Write([]byte(p.To + "\n"))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *Ping) UnmarshalBinary(data []byte) error {
	buf := bufio.NewReader(bytes.NewReader(data))
	var err error
	p.From, err = buf.ReadString('\n')
	if err != nil {
		return err
	}
	p.To, err = buf.ReadString('\n')
	if err != nil {
		return err
	}
	return nil
}
