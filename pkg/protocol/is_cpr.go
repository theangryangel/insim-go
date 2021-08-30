package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspCpr ...
const (
	IspCpr = 20
)

// Cpr ...
type Cpr struct {
	ReqI uint8 `struct:"uint8"`
	Ucid uint8 `struct:"uint8"`

	PName string `struct:"[24]byte"`
	Plate string `struct:"[8]byte"`
}

// UnmarshalInsim ...
func (p *Cpr) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Cpr) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Cpr) Type() uint8 {
	return IspCpr
}

// NewCpr ...
func NewCpr() Packet {
	return &Cpr{}
}

// New ...
func (p *Cpr) New() Packet {
	return NewCpr()
}
