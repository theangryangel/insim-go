package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspPla ...
const (
	IspPla = 28

	PitlaneFactExit      = 0
	PitlaneFactEnter     = 1
	PitlaneFactNoPurpose = 2
	PitlaneFactDT        = 3
	PitlaneFactSG        = 4
)

// Pla ...
type Pla struct {
	ReqI uint8 `struct:"uint8"`
	Plid uint8 `struct:"uint8"`
	Fact uint8 `struct:"uint8"`
	Sp1  uint8 `struct:"uint8"`
	Sp2  uint8 `struct:"uint8"`
	Sp3  uint8 `struct:"uint8"`
}

// Entering ...
func (p *Pla) Entering() bool {
	switch p.Fact {
	case PitlaneFactEnter, PitlaneFactNoPurpose, PitlaneFactDT, PitlaneFactSG:
		return true
	default:
		return false
	}
}

// Exiting ...
func (p *Pla) Exiting() bool {
	return p.Fact == PitlaneFactExit
}

// NoPurpose ...
func (p *Pla) NoPurpose() bool {
	return p.Fact == PitlaneFactNoPurpose
}

// StopGo ...
func (p *Pla) StopGo() bool {
	return p.Fact == PitlaneFactSG
}

// DriveThrough ...
func (p *Pla) DriveThrough() bool {
	return p.Fact == PitlaneFactDT
}

// UnmarshalInsim ...
func (p *Pla) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Pla) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Pla) Type() uint8 {
	return IspPla
}

// NewPla ...
func NewPla() Packet {
	return &Pla{}
}

// New ...
func (p *Pla) New() Packet {
	return NewPla()
}
