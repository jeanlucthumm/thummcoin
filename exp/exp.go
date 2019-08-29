package main

import (
	"log"
	"net"
	"time"
)

func client() {
	var conn net.Conn
	var err error
	// retry server until it is up
	for {
		conn, err = net.Dial("tcp", ":8081")
		if err == nil {
			break
		}
		log.Println(err)
		time.Sleep(time.Second)
	}
	_, err = conn.Write([]byte("request"))
	if err != nil {
		log.Println(err)
		return
	}

	var buf []byte
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("From server: %s\n", buf[:n])
}

func server() {
	ln, _ := net.Listen("tcp", ":8081")
	for {
		conn, _ := ln.Accept()
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	var buf []byte
	n, err := conn.Read(buf)
	if err != nil {
		return
	}
	log.Printf("Server got: %s\n", buf)

	if string(buf[:n]) == "request" {
		_, _ = conn.Write([]byte("response"))
	}
}

func main() {
	go client()
	server()
}
