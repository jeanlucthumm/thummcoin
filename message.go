package main

import (
	"bytes"
	"github.com/lunixbochs/struc"
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
	err := struc.Pack(buf, t)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (t *Transaction) UnmarshalBinary(data []byte) error {
	buf := bytes.NewReader(data)
	return struc.Unpack(buf, t)
}
