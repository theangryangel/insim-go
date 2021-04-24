package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_VER = 2
)

type Ver struct {
	ReqI uint8 `struct:"uint8"`
	Zero uint8 `struct:"uint8"`

	Version  string `struct:"[8]byte"`
	Product  string `struct:"[6]byte"`
	InSimVer uint8  `struct:"uint8"`
	Spare    uint8  `struct:"uint8"`
}

func (p *Ver) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Ver) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Ver) Type() (id uint8) {
	return ISP_VER
}

func NewVer() Packet {
	return &Ver{}
}

func (p *Ver) New() Packet {
	return NewVer()
}
