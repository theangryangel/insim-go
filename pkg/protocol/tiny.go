package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_TINY = 3 //  - both ways		: multi purpose
)

type Tiny struct {
	ReqI uint8 `struct:"uint8"`
	SubT uint8 `struct:"uint8"`
}

func (p *Tiny) Unmarshal(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Tiny) Marshal() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Tiny) Type() (id uint8) {
	return ISP_TINY
}

func NewTiny() (Packet) {
	return &Tiny{
	}
}

func (p *Tiny) New() (Packet) {
	return NewTiny()
}
