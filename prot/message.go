package prot

const (
	PING  = iota
	PLIST
)

type Message struct {
	ID   int
	data []byte
}
