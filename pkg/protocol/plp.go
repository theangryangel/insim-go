package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_PLP = 22
)

type Plp struct {
	ReqI uint8 `struct:"uint8"`
	Plid uint8 `struct:"uint8"`
}

func (p *Plp) Unmarshal(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Plp) Marshal() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Plp) Type() (uint8) {
	return ISP_PLL
}

func NewPlp() (Packet) {
	return &Plp{
	}
}

func (p *Plp) New() (Packet) {
	return NewPlp()
}
