package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspVtn ...
const (
	IspVtn = 16

	VoteNone    = 0
	VoteEnd     = 1
	VoteRestart = 2
	VoteQualify = 3
)

// Vtn ...
type Vtn struct {
	ReqI uint8 `struct:"uint8"`
	Zero uint8 `struct:"uint8"`

	Ucid   uint8 `struct:"uint8"`
	Action uint8 `struct:"uint8"`
	Sp2    uint8 `struct:"uint8"`
	Sp3    uint8 `struct:"uint8"`
}

// Voting ...
func (p *Vtn) Voting() bool {
	return p.Action != VoteNone
}

// Restart ...
func (p *Vtn) Restart() bool {
	return p.Action == VoteRestart
}

// End ...
func (p *Vtn) End() bool {
	return p.Action == VoteEnd
}

// Qualify ...
func (p *Vtn) Qualify() bool {
	return p.Action == VoteQualify
}

// UnmarshalInsim ...
func (p *Vtn) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Vtn) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Vtn) Type() uint8 {
	return IspVtn
}

// NewVtn ...
func NewVtn() Packet {
	return &Vtn{}
}

// New ...
func (p *Vtn) New() Packet {
	return NewVtn()
}
