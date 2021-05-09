package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_REO = 36
)

type Reo struct {
	ReqI uint8 `struct:"uint8"`

	NumP uint8     `struct:"uint8"`
	Plid [40]uint8 `struct:"[40]uint8"`
}

func (p *Reo) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Reo) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Reo) Type() uint8 {
	return ISP_REO
}

func NewReo() Packet {
	return &Reo{}
}

func (p *Reo) New() Packet {
	return NewReo()
}
