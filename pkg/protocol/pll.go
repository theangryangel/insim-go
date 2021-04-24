package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_PLL = 23
)

type Pll struct {
	ReqI uint8 `struct:"uint8"`
	Plid uint8 `struct:"uint8"`
}

func (p *Pll) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Pll) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Pll) Type() uint8 {
	return ISP_PLL
}

func NewPll() Packet {
	return &Pll{}
}

func (p *Pll) New() Packet {
	return NewPll()
}
