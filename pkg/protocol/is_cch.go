package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspCch ...
const (
	// IspCch is the packet ID for a CCH packet
	IspCch = 29
)

// Cch is a Camera Change Packet
type Cch struct {
	ReqI uint8 `struct:"uint8"`

	Plid   uint8 `struct:"uint8"`
	Camera uint8 `struct:"uint8"`
	Sp1    uint8 `struct:"uint8"`
	Sp2    uint8 `struct:"uint8"`
	Sp3    uint8 `struct:"uint8"`
}

// UnmarshalInsim unpacks a Cnl Packet
func (p *Cch) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim packs a Cnl Packet
func (p *Cch) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type returns the packet ID code
func (p *Cch) Type() uint8 {
	return IspCch
}

// NewCch returns a new Cch
func NewCch() Packet {
	return &Cch{}
}

// New returns a new Cch
func (p *Cch) New() Packet {
	return NewCch()
}
