package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspMci ...
const (
	IspMci = 38
)

// MciInfo ...
type MciInfo struct {
	Node      uint16 `struct:"uint16"`
	Lap       uint16 `struct:"uint16"`
	Plid      uint8  `struct:"uint8"`
	Position  uint8  `struct:"uint8"`
	Info      uint8  `struct:"uint8"`
	Sp3       uint8  `struct:"uint8"`
	X         int32  `struct:"int32"`
	Y         int32  `struct:"int32"`
	Z         int32  `struct:"int32"`
	Speed     uint16 `struct:"uint16"`
	Direction uint16 `struct:"uint16"`
	Heading   uint16 `struct:"uint16"`
	AngVel    int16  `struct:"int16"`
}

// Kmph ...
func (p *MciInfo) Kmph() float32 {
	// raw speed to Kilometers per hour
	return float32(p.Speed) / 91.02
}

// Mph ...
func (p *MciInfo) Mph() float32 {
	// raw speed to Miles Per Hour
	return float32(p.Speed) / 146.486067
}

// Mps ...
func (p *MciInfo) Mps() float32 {
	// raw speed to Meters per second
	return float32(p.Speed) / 32768
}

// Mci ...
type Mci struct {
	ReqI uint8 `struct:"uint8"`
	NumC uint8 `struct:"uint8,sizeof=Info"`
	Info []MciInfo
}

// UnmarshalInsim ...
func (p *Mci) UnmarshalInsim(data []byte) (err error) {
	err = restruct.Unpack(data, binary.LittleEndian, p)
	if err != nil {
		return err
	}

	return nil
}

// MarshalInsim ...
func (p *Mci) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Mci) Type() uint8 {
	return IspMci
}

// NewMci ...
func NewMci() Packet {
	return &Mci{}
}

// New ...
func (p *Mci) New() Packet {
	return NewMci()
}
