package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"time"
)

// IspSpx ...
const (
	IspSpx = 25
)

// Spx ...
type Spx struct {
	ReqI uint8 `struct:"uint8"`

	Plid uint8 `struct:"uint8"`

	STime uint32 `struct:"uint32"`
	ETime uint32 `struct:"uint32"`

	Split    uint8 `struct:"uint8"`
	Penalty  uint8 `struct:"uint8"`
	NumStops uint8 `struct:"uint8"`
	Fuel200  uint8 `struct:"uint8"`
}

// Fuel ...
func (p *Spx) Fuel() uint8 {
	// TODO magic decoding
	return p.Fuel200
}

// SplitTime ...
func (p *Spx) SplitTime() time.Duration {
	return (time.Duration(p.STime) * time.Millisecond)
}

// ElapsedTime ...
func (p *Spx) ElapsedTime() time.Duration {
	return (time.Duration(p.ETime) * time.Millisecond)
}

// UnmarshalInsim ...
func (p *Spx) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Spx) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Spx) Type() (id uint8) {
	return IspSpx
}

// NewSpx ...
func NewSpx() Packet {
	return &Spx{}
}

// New ...
func (p *Spx) New() Packet {
	return NewSpx()
}
