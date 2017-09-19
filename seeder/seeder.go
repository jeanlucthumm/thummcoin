package seeder

import (
	"net"
	"log"
)

var ips = []string {
	":8080",
	":8081",
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

}
