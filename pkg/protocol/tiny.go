package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_TINY = 3 //  - both ways		: multi purpose

	TINY_NONE  = 0
	TINY_VER   = 1
	TINY_CLOSE = 2
	TINY_PING  = 3
	TINY_REPLY = 4
	TINY_VTC   = 5
	TINY_SCP   = 6
	TINY_SST   = 7
	TINY_GTH   = 8
	TINY_MPE   = 9
	TINY_ISM   = 10
	TINY_REN   = 11
	TINY_CLR   = 12
	TINY_NCN   = 13
	TINY_NPL   = 14
	TINY_RES   = 15
	TINY_NLP   = 16
	TINY_MCI   = 17
	TINY_REO   = 18
	TINY_RST   = 19
	TINY_AXI   = 20
	TINY_AXC   = 21
	TINY_RIP   = 22
	TINY_NCI   = 23
	TINY_ALC   = 24
	TINY_AXM   = 25
	TINY_SLC   = 26
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

func NewTiny() Packet {
	return &Tiny{}
}

func (p *Tiny) New() Packet {
	return NewTiny()
}
