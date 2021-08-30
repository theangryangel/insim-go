package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspVer ...
const (
	IspVer = 2
)

// Ver ...
type Ver struct {
	ReqI uint8 `struct:"uint8"`
	Zero uint8 `struct:"uint8"`

	Version  string `struct:"[8]byte"`
	Product  string `struct:"[6]byte"`
	InSimVer uint8  `struct:"uint8"`
	Spare    uint8  `struct:"uint8"`
}

// UnmarshalInsim ...
func (p *Ver) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Ver) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Ver) Type() (id uint8) {
	return IspVer
}

// NewVer ...
func NewVer() Packet {
	return &Ver{}
}

// New ...
func (p *Ver) New() Packet {
	return NewVer()
}
