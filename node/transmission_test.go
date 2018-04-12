package node

import (
	"testing"
	"net"
	"github.com/jeanlucthumm/thummcoin/prot"
	"fmt"
)

const (
	testingAddr = ":9090"
	rwtimeout   = 2
)

// newServer creates a server listening on addr which processes connection through handler.
// The server will send true on started once it is ready to accept connections and will
// terminate execution if anything is received on done.
func newServer(addr string, started, done chan bool, handler func(net.Conn) error) {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Errorf("server could not listen")
		return
	}
	defer ln.Close()
	started <- true

	for {
		select {
		case <-done:
			return
		default:
			conn, err := ln.Accept()
			if err != nil {
				fmt.Errorf("server could not accept")
				return
			}

			err = handler(conn)
			if err != nil {
				fmt.Errorf("server got error from handler: %v", err)
				return
			}

			conn.Close()
		}
	}
}

func TestAwaitSend(t *testing.T) {
	started := make(chan bool, 1)
	done := make(chan bool, 1)

	// server simply returns everything it's sent
	go newServer(testingAddr, started, done, func(conn net.Conn) error {
		for {
			select {
			case <-done:
				return nil
			default:
				// read data and reply with the same thing
				buf := make([]byte, 4096)
				n, err := conn.Read(buf)
				if err != nil {
					return err
				}

				_, err = conn.Write(buf[:n])
				if err != nil {
					return err
				}
			}
		}
	})

	<-started
	m, err := prot.MakePingMessage(&prot.Ping{From: "me", To: "you"})
	conn, err := net.Dial("tcp", testingAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	reply, err := sendAwait(conn, m)
	if err != nil {
		t.Error(err)
		return
	}

	if !reply.Equal(m) {
		t.Fail()
	}
	done <- true
}

func TestPingExchange(t *testing.T) {
	started := make(chan bool, 1)
	done := make(chan bool, 1)

	// server decodes ping and sends another valid one
	go newServer(testingAddr, started, done, func(conn net.Conn) error {
		// extract ping data
		m, err := prot.Receive(conn)
		if err != nil {
			return err
		}

		pm, err := prot.DecodeMessage(m)
		if err != nil {
			return err
		}

		p, ok := pm.(*prot.Ping)
		if !ok {
			t.Log("server did not recieve proper ping data")
			t.Fail()
			return nil
		}

		// send reply
		preply, err := prot.MakePingMessage(p)
		if err != nil {
			return err
		}

		err = prot.Send(conn, preply)
		if err != nil {
			return err
		}

		return nil
	})

	<-started
	conn, err := net.Dial("tcp", testingAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	err = pingExchange(conn, &prot.Ping{From: "me", To: "you"})
	if err != nil {
		t.Error(err)
	}
}
