package main

import (
	"bytes"
	"net"
	"encoding/binary"
	"fmt"
	"bufio"
	"strings"
)

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
	num  uint32
	list []net.Addr
}

func (p *PeerList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, p.num)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, addr := range p.list {
		buf.Write([]byte(addr.String() + "\n"))
	}
	return buf.Bytes(), nil
}

func (p *PeerList) UnmarshalBinary(data []byte) error {
	r := bytes.NewReader(data)
	err := binary.Read(r, binary.BigEndian, &p.num)
	if err != nil {
		fmt.Println(err)
		return err
	}

	buf := bufio.NewReader(r)
	for i := uint32(0); i < p.num; i++ {
		str, err := buf.ReadString('\n')
		if err != nil {
			return err
		}
		addr, err := net.ResolveTCPAddr("tcp", strings.Trim(str, "\n"))
		if err != nil {
			return err
		}
		p.list = append(p.list, addr)
	}
	return nil
}
