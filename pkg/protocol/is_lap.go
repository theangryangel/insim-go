package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"time"
)

// IspLap ...
const (
	IspLap = 24
)

// Lap ...
type Lap struct {
	ReqI uint8 `struct:"uint8"`
	Plid uint8 `struct:"uint8"`

	LTime uint32 `struct:"uint32"`
	ETime uint32 `struct:"uint32"`

	LapsDone uint16 `struct:"uint16"`
	Flags    uint16 `struct:"uint16"`

	Sp0      uint8 `struct:"uint8"`
	Penalty  uint8 `struct:"uint8"`
	NumStops uint8 `struct:"uint8"`
	Fuel200  uint8 `struct:"uint8"`
}

// Fuel ...
func (p *Lap) Fuel() uint8 {
	// TODO magic decoding
	return p.Fuel200
}

// LapTime ...
func (p *Lap) LapTime() time.Duration {
	return (time.Duration(p.LTime) * time.Millisecond)
}

// ElapsedTime ...
func (p *Lap) ElapsedTime() time.Duration {
	return (time.Duration(p.ETime) * time.Millisecond)
}

// UnmarshalInsim ...
func (p *Lap) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Lap) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Lap) Type() (id uint8) {
	return IspLap
}

// NewLap ...
func NewLap() Packet {
	return &Lap{}
}

// New ...
func (p *Lap) New() Packet {
	return NewLap()
}
