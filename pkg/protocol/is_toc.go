package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_TOC = 31
)

type Toc struct {
	ReqI uint8 `struct:"uint8"`
	Plid uint8 `struct:"uint8"`

	OldUcid uint8 `struct:"uint8"`
	NewUcid uint8 `struct:"uint8"`
	Sp2     uint8 `struct:"uint8"`
	Sp3     uint8 `struct:"uint8"`
}

func (p *Toc) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Toc) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Toc) Type() (id uint8) {
	return ISP_TOC
}

func NewToc() Packet {
	return &Toc{}
}

func (p *Toc) New() Packet {
	return NewToc()
}
