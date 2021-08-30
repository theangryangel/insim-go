package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspPen ...
const (
	IspPen = 30

	// Penalty types
	PenaltyNone    = 0
	PenaltyDT      = 1
	PenaltyDTValid = 2
	PenaltySG      = 3
	PenaltySGValid = 4
	Penalty30sec   = 5
	Penalty45sec   = 6

	// Penatly Reasons
	PenaltyReasonUnknown    = 0
	PenaltyReasonAdmin      = 1
	PenaltyReasonWrongWay   = 2
	PenaltyReasonFalseStart = 3
	PenaltyReasonSpeeding   = 4
	PenaltyReasonStopShort  = 5
	PenaltyReasonStopLate   = 6
)

// Pen ...
type Pen struct {
	ReqI uint8 `struct:"uint8"`
	Plid uint8 `struct:"uint8"`

	OldPen uint8 `struct:"uint8"`
	NewPen uint8 `struct:"uint8"`
	Reason uint8 `struct:"uint8"`
	Sp3    uint8 `struct:"uint8"`
}

// UnmarshalInsim ...
func (p *Pen) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Pen) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Pen) Type() uint8 {
	return IspPen
}

// NewPen ...
func NewPen() Packet {
	return &Pen{}
}

// New ...
func (p *Pen) New() Packet {
	return NewPen()
}
