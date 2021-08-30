package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspCon ...
const (
	// IspCon is the packet ID for a CON packet
	IspCon = 50
)

// CarContact describes a vehicle in a contact
type CarContact struct {
	Plid uint8 `struct:"uint8"`
	Info uint8 `struct:"uint8"`
	Sp2  uint8 `struct:"uint8"`

	Steer     int8  `struct:"int8"`
	ThrBrk    uint8 `struct:"uint8"`
	CluHan    uint8 `struct:"uint8"`
	GearSp    uint8 `struct:"uint8"`
	Speed     uint8 `struct:"uint8"`
	Direction uint8 `struct:"uint8"`
	Heading   uint8 `struct:"uint8"`

	AccelF int8 `struct:"int8"`
	AccelR int8 `struct:"int8"`

	X int16 `struct:"int16"`
	Y int16 `struct:"int16"`
}

// Con describes the contact between 2 players
type Con struct {
	ReqI uint8 `struct:"uint8"`
	Zero uint8 `struct:"uint8"`

	SpClose uint16 `struct:"uint16"`
	Time    uint16 `struct:"uint16"`

	A CarContact
	B CarContact
}

// UnmarshalInsim unpacks a packet
func (p *Con) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim packs a packet
func (p *Con) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type returns the packet ID
func (p *Con) Type() (id uint8) {
	return IspCon
}

// NewCon creates a new Con
func NewCon() Packet {
	return &Con{}
}

// New is a proxy for NewCon
func (p *Con) New() Packet {
	return NewCon()
}
