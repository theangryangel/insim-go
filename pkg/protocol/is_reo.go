package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspReo ...
const (
	IspReo = 36
)

// Reo ...
type Reo struct {
	ReqI uint8 `struct:"uint8"`

	NumP uint8     `struct:"uint8"`
	Plid [40]uint8 `struct:"[40]uint8"`
}

// UnmarshalInsim ...
func (p *Reo) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Reo) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Reo) Type() uint8 {
	return IspReo
}

// NewReo ...
func NewReo() Packet {
	return &Reo{}
}

// New ...
func (p *Reo) New() Packet {
	return NewReo()
}
