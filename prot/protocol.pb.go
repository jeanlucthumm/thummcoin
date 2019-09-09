// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protocol.proto

/*
Package prot is a generated protocol buffer package.

It is generated from these files:
	protocol.proto

It has these top-level messages:
	Message
	Request
	PeerList
*/
package prot

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Type int32

const (
	Type_PING      Type = 0
	Type_REQ       Type = 1
	Type_PEER_LIST Type = 2
	Type_TEXT      Type = 3
)

var Type_name = map[int32]string{
	0: "PING",
	1: "REQ",
	2: "PEER_LIST",
	3: "TEXT",
}
var Type_value = map[string]int32{
	"PING":      0,
	"REQ":       1,
	"PEER_LIST": 2,
	"TEXT":      3,
}

func (x Type) String() string {
	return proto.EnumName(Type_name, int32(x))
}
func (Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Request_Type int32

const (
	Request_PEER_LIST Request_Type = 0
	Request_IP_SELF   Request_Type = 1
)

var Request_Type_name = map[int32]string{
	0: "PEER_LIST",
	1: "IP_SELF",
}
var Request_Type_value = map[string]int32{
	"PEER_LIST": 0,
	"IP_SELF":   1,
}

func (x Request_Type) String() string {
	return proto.EnumName(Request_Type_name, int32(x))
}
func (Request_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

type Message struct {
	Type Type   `protobuf:"varint,1,opt,name=type,enum=prot.Type" json:"type,omitempty"`
	From string `protobuf:"bytes,2,opt,name=from" json:"from,omitempty"`
	To   string `protobuf:"bytes,3,opt,name=to" json:"to,omitempty"`
	// 2 - 10 reserved for inner protocol
	Data []byte `protobuf:"bytes,11,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message) GetType() Type {
	if m != nil {
		return m.Type
	}
	return Type_PING
}

func (m *Message) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *Message) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *Message) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type Request struct {
	Type Request_Type `protobuf:"varint,1,opt,name=type,enum=prot.Request_Type" json:"type,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Request) GetType() Request_Type {
	if m != nil {
		return m.Type
	}
	return Request_PEER_LIST
}

type PeerList struct {
	Peers []*PeerList_Peer `protobuf:"bytes,1,rep,name=peers" json:"peers,omitempty"`
}

func (m *PeerList) Reset()                    { *m = PeerList{} }
func (m *PeerList) String() string            { return proto.CompactTextString(m) }
func (*PeerList) ProtoMessage()               {}
func (*PeerList) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *PeerList) GetPeers() []*PeerList_Peer {
	if m != nil {
		return m.Peers
	}
	return nil
}

type PeerList_Peer struct {
	Address string `protobuf:"bytes,1,opt,name=address" json:"address,omitempty"`
}

func (m *PeerList_Peer) Reset()                    { *m = PeerList_Peer{} }
func (m *PeerList_Peer) String() string            { return proto.CompactTextString(m) }
func (*PeerList_Peer) ProtoMessage()               {}
func (*PeerList_Peer) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }

func (m *PeerList_Peer) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func init() {
	proto.RegisterType((*Message)(nil), "prot.Message")
	proto.RegisterType((*Request)(nil), "prot.Request")
	proto.RegisterType((*PeerList)(nil), "prot.PeerList")
	proto.RegisterType((*PeerList_Peer)(nil), "prot.PeerList.Peer")
	proto.RegisterEnum("prot.Type", Type_name, Type_value)
	proto.RegisterEnum("prot.Request_Type", Request_Type_name, Request_Type_value)
}

func init() { proto.RegisterFile("protocol.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 263 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0x41, 0x4b, 0xc4, 0x30,
	0x10, 0x85, 0x37, 0x6d, 0xb4, 0xdb, 0xa9, 0x96, 0x32, 0x5e, 0x82, 0x07, 0x29, 0x3d, 0x48, 0xf5,
	0xb0, 0x87, 0xfa, 0x1b, 0xaa, 0x14, 0xaa, 0xd4, 0x6c, 0x45, 0x6f, 0x4b, 0xb5, 0xa3, 0x08, 0x4a,
	0x6a, 0x12, 0x0f, 0xfb, 0xef, 0x25, 0xe9, 0x2e, 0xb8, 0xa7, 0xbc, 0x37, 0xf3, 0xc8, 0xf7, 0x18,
	0x48, 0x27, 0xad, 0xac, 0x7a, 0x53, 0x5f, 0x2b, 0x2f, 0x90, 0xbb, 0xa7, 0x18, 0x20, 0xba, 0x27,
	0x63, 0x86, 0x0f, 0xc2, 0x0b, 0xe0, 0x76, 0x3b, 0x91, 0x60, 0x39, 0x2b, 0xd3, 0x0a, 0x7c, 0x6c,
	0xd5, 0x6f, 0x27, 0x92, 0x7e, 0x8e, 0x08, 0xfc, 0x5d, 0xab, 0x6f, 0x11, 0xe4, 0xac, 0x8c, 0xa5,
	0xd7, 0x98, 0x42, 0x60, 0x95, 0x08, 0xfd, 0x24, 0xb0, 0xca, 0x65, 0xc6, 0xc1, 0x0e, 0x22, 0xc9,
	0x59, 0x79, 0x22, 0xbd, 0x2e, 0x9e, 0x20, 0x92, 0xf4, 0xf3, 0x4b, 0xc6, 0xe2, 0xe5, 0x01, 0x02,
	0x67, 0xc4, 0x6e, 0xf9, 0x0f, 0x55, 0x14, 0xc0, 0x9d, 0xc3, 0x53, 0x88, 0xbb, 0xba, 0x96, 0x9b,
	0xb6, 0x59, 0xf7, 0xd9, 0x02, 0x13, 0x88, 0x9a, 0x6e, 0xb3, 0xae, 0xdb, 0xdb, 0x8c, 0x15, 0xcf,
	0xb0, 0xec, 0x88, 0x74, 0xfb, 0x69, 0x2c, 0x5e, 0xc1, 0xd1, 0x44, 0xa4, 0x8d, 0x60, 0x79, 0x58,
	0x26, 0xd5, 0xd9, 0xfc, 0xf1, 0x7e, 0xed, 0x85, 0x9c, 0x13, 0xe7, 0x39, 0x70, 0x67, 0x51, 0x40,
	0x34, 0x8c, 0xa3, 0x26, 0x63, 0x7c, 0x9b, 0x58, 0xee, 0xed, 0x75, 0xb5, 0x83, 0x2f, 0x81, 0x77,
	0xcd, 0xc3, 0x5d, 0xb6, 0xc0, 0x08, 0x42, 0x59, 0x3f, 0x66, 0xec, 0xb0, 0x4f, 0xe0, 0x12, 0x7d,
	0xfd, 0xd2, 0x67, 0xe1, 0xeb, 0xb1, 0xbf, 0xe9, 0xcd, 0x5f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc4,
	0xf5, 0x9f, 0xae, 0x65, 0x01, 0x00, 0x00,
}
