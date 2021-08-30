package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspPlp ...
const (
	IspPlp = 22
)

// Plp ...
type Plp struct {
	ReqI uint8 `struct:"uint8"`
	Plid uint8 `struct:"uint8"`
}

// UnmarshalInsim ...
func (p *Plp) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Plp) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Plp) Type() uint8 {
	return IspPlp
}

// NewPlp ...
func NewPlp() Packet {
	return &Plp{}
}

// New ...
func (p *Plp) New() Packet {
	return NewPlp()
}
