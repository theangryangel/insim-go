package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_CCH = 29
)

type Cch struct {
	ReqI uint8 `struct:"uint8"`

	Plid   uint8 `struct:"uint8"`
	Camera uint8 `struct:"uint8"`
	Sp1    uint8 `struct:"uint8"`
	Sp2    uint8 `struct:"uint8"`
	Sp3    uint8 `struct:"uint8"`
}

func (p *Cch) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Cch) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Cch) Type() uint8 {
	return ISP_CCH
}

func NewCch() Packet {
	return &Cch{}
}

func (p *Cch) New() Packet {
	return NewCch()
}
