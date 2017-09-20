package seeder

import (
	"net"
	"log"
	"github.com/jeanlucthumm/thummcoin/prot"
)

// SeederIPs holds the endpoints addresses of the seeders
var SeederIPs = []string{
	":8090",
}

var ips = []string{
	":8080",
	":8081",
	":8082",
}

var peerList prot.PeerList
var peerBinary []byte

func init() {
	for _, ip := range ips {
		addr, err := net.ResolveTCPAddr("tcp", ip)
		if err != nil {
			log.Println("Could not reolve tcp addr:", ip)
			continue
		}
		peerList.List = append(peerList.List, addr)
		peerList.Num++
	}

	var err error
	peerBinary, err = peerList.MarshalBinary()
	if err != nil {
		log.Println("Could not marshal binary of peerList")
	}
}

func Start(network string, loc string) error {
	addr, err := net.ResolveTCPAddr(network, loc)
	if err != nil {
		log.Fatalln("Could not start seeder:", err)
		return err
	}

	ln, err := net.Listen(addr.Network(), addr.String())
	if err != nil {
		log.Fatalln("Could not start seeder:", err)
		return err
	}
	log.Println("Seeder listening on:", addr.String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		log.Println("New connection:", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	conn.Write(peerBinary)
}
