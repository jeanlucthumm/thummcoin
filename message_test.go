package main

import (
	"testing"
	"fmt"
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
