package seeder

import (
	"testing"
	"fmt"
	"net"
	"github.com/jeanlucthumm/thummcoin/node"
)

func TestStart(t *testing.T) {
	go Start("tcp", ":8090")
	conn, err := net.Dial("tcp", ":8090")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil || n <= 0 {
		fmt.Println(err)
		t.Fail()
	}

	var plist node.PeerList
	err = plist.UnmarshalBinary(buf)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
