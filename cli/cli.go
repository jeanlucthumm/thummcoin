package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/jeanlucthumm/thummcoin/p2p"
	"github.com/jeanlucthumm/thummcoin/prot"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io"
)

var log = logrus.WithField("mod", "cli")

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
			log.Errorf("Failed to parse command: %s\n", err)
		}
		rest := line[len(cmd):]
		cmd = cmd[:len(cmd)-1] // remove included delim
		switch cmd {
		case "msg":
			if err = msg(reader, node); err != nil {
				log.Errorf("Failed to process msg command: %s\n", err)
			}
		case "peers":
			if err = peers(rest, node); err != nil {
				log.Errorf("Failed to process peers command: %s\n", err)
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

func peers(line []byte, node *p2p.Node) error {
	scan := bufio.NewScanner(bytes.NewReader(line))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("incomplete peers command")
		return nil
	}
	switch scan.Text() {
	case "list":
		ips := node.ListPeers()
		for _, ip := range ips {
			fmt.Println(ip.String())
		}
	default:
		fmt.Println("unknown peers command")
		return nil
	}
	return nil
}
