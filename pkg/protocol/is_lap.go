package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"time"
)

const (
	ISP_LAP = 24
)

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

func (p *Lap) Fuel() uint8 {
	// TODO magic decoding
	return p.Fuel200
}

func (p *Lap) LapTime() time.Duration {
	return (time.Duration(p.LTime) * time.Millisecond)
}

func (p *Lap) ElapsedTime() time.Duration {
	return (time.Duration(p.ETime) * time.Millisecond)
}

func (p *Lap) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Lap) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Lap) Type() (id uint8) {
	return ISP_LAP
}

func NewLap() Packet {
	return &Lap{}
}

func (p *Lap) New() Packet {
	return NewLap()
}
