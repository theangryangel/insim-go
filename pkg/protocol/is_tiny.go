package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspTiny ...
const (
	IspTiny = 3 //  - both ways		: multi purpose

	TinyNone  = 0
	TinyVer   = 1
	TinyClose = 2
	TinyPing  = 3
	TinyReply = 4
	TinyVTC   = 5
	TinySCP   = 6
	TinySST   = 7
	TinyGTH   = 8
	TinyMPE   = 9
	TinyISM   = 10
	TinyREN   = 11
	TinyCLR   = 12
	TinyNCN   = 13
	TinyNPL   = 14
	TinyRES   = 15
	TinyNLP   = 16
	TinyMCI   = 17
	TinyREO   = 18
	TinyRST   = 19
	TinyAXI   = 20
	TinyAXC   = 21
	TinyRIP   = 22
	TinyNCI   = 23
	TinyALC   = 24
	TinyAXM   = 25
	TinySLC   = 26
)

// Tiny ...
type Tiny struct {
	ReqI uint8 `struct:"uint8"`
	SubT uint8 `struct:"uint8"`
}

// UnmarshalInsim ...
func (p *Tiny) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Tiny) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Tiny) Type() (id uint8) {
	return IspTiny
}

// NewTiny ...
func NewTiny() Packet {
	return &Tiny{}
}

// New ...
func (p *Tiny) New() Packet {
	return NewTiny()
}
