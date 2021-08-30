package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspToc ...
const (
	IspToc = 31
)

// Toc ...
type Toc struct {
	ReqI uint8 `struct:"uint8"`
	Plid uint8 `struct:"uint8"`

	OldUcid uint8 `struct:"uint8"`
	NewUcid uint8 `struct:"uint8"`
	Sp2     uint8 `struct:"uint8"`
	Sp3     uint8 `struct:"uint8"`
}

// UnmarshalInsim ...
func (p *Toc) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Toc) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Toc) Type() (id uint8) {
	return IspToc
}

// NewToc ...
func NewToc() Packet {
	return &Toc{}
}

// New ...
func (p *Toc) New() Packet {
	return NewToc()
}
