package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_MCI = 38
)

type MciInfo struct {
	Node uint16 `struct:"uint16"`
	Lap uint16 `struct:"uint16"`
	Plid uint8 `struct:"uint8"`
	Position uint8 `struct:"uint8"`
	Info uint8 `struct:"uint8"`
	Sp3 uint8 `struct:"uint8"`
	X int32 `struct:"int32"`
	Y int32 `struct:"int32"`
	Z int32 `struct:"int32"`
	Speed uint16 `struct:"uint16"`
	Direction uint16 `struct:"uint16"`
	Heading uint16 `struct:"uint16"`
	AngVel int16 `struct:"int16"`
}

func (p *MciInfo) Kmph() (float32) {
	// raw speed to Kilometers per hour
	return float32(p.Speed) / 91.02
}

func (p *MciInfo) Mph() (float32) {
	// raw speed to Miles Per Hour
	return float32(p.Speed) / 146.486067
}

func (p *MciInfo) Mps() (float32) {
	// raw speed to Meters per second
	return float32(p.Speed) / 32768
}

type Mci struct {
	ReqI uint8 `struct:"uint8"`
	NumC uint8 `struct:"uint8,sizeof=Info"`
	Info []MciInfo
}

func (p *Mci) Unmarshal(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Mci) Marshal() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Mci) Type() (uint8) {
	return ISP_MCI
}

func NewMci() (Packet) {
	return &Mci{
	}
}

func (p *Mci) New() (Packet) {
	return NewMci()
}
