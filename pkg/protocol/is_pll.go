package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspPll ...
const (
	IspPll = 23
)

// Pll ...
type Pll struct {
	ReqI uint8 `struct:"uint8"`
	Plid uint8 `struct:"uint8"`
}

// UnmarshalInsim ...
func (p *Pll) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Pll) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Pll) Type() uint8 {
	return IspPll
}

// NewPll ...
func NewPll() Packet {
	return &Pll{}
}

// New ...
func (p *Pll) New() Packet {
	return NewPll()
}
