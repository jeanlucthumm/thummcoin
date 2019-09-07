package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jeanlucthumm/thummcoin/p2p"
	"github.com/jeanlucthumm/thummcoin/prot"
	"github.com/pkg/errors"
	"io"
	"log"
)

func Interpret(in io.Reader, node *p2p.Node) {
	lScan := bufio.NewScanner(in)
	lScan.Split(bufio.ScanLines)
	for lScan.Scan() {
		line := append(lScan.Bytes(), '\n') // for helper funcs
		reader := bufio.NewReader(bytes.NewReader(line))

		cmd, err := reader.ReadString(' ')
		if err == io.EOF {
			fmt.Println("incomplete command")
		} else if err != nil {
			log.Printf("Failed to parse command: %s\n", err)
		}
		cmd = cmd[:len(cmd)-1] // remove included delim

		switch cmd {
		case "msg":
			if err = msg(reader, node); err != nil {
				log.Printf("Failed to process msg command: %s\n", err)
			}
		}
	}
}

func msg(reader *bufio.Reader, node *p2p.Node) error {
	w, err := reader.ReadString(' ')
	if err == io.EOF {
		fmt.Println("incomplete command")
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "parse token")
	}
	w = w[:len(w)-1]

	switch w {
	case "broadcast":
		msg, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Println("incomplete command")
			return nil
		}
		if err != nil {
			return errors.Wrap(err, "parse token")
		}
		msg = msg[:len(msg)-1]
		broadcast(msg, node)
		return nil
	default:
		fmt.Println("unknown command")
		return nil
	}
}

func broadcast(msg string, node *p2p.Node) {
	m := &p2p.Message{
		Kind: prot.Type_TEXT,
		Data: []byte(msg),
	}
	node.Broadcast <- m
}
