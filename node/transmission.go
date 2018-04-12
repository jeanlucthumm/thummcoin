package node

import (
	"net"
	"github.com/jeanlucthumm/thummcoin/prot"
	"time"
)

// pingExchange sends a ping message along conn and waits for a ping reply.
func pingExchange(conn net.Conn, p *prot.Ping) error {
	// construct message and wait for reply
	chal, err := prot.MakePingMessage(p)
	if err != nil {
		return err
	}
	rm, err := sendAwait(conn, chal)
	if err != nil {
		return err
	}

	// check reply
	rt, err := prot.DecodeMessage(rm)
	if err != nil {
		return err
	}
	if r, ok := rt.(*prot.Ping); ok {
		if r.From == p.To && r.To == p.From {
			return nil
		}
	}
	return err
}

// sendAwait sends a message along conn and returns the reply
func sendAwait(conn net.Conn, m *prot.Message) (*prot.Message, error) {
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	if err := prot.Send(conn, m); err != nil {
		return nil, err
	}
	conn.SetDeadline(time.Now().Add(2 * time.Second))
	reply, err := prot.Receive(conn)
	return reply, err
}
