package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_PEN = 30

	// Penalty types
	PENALTY_NONE     = 0
	PENALTY_DT       = 1
	PENALTY_DT_VALID = 2
	PENALTY_SG       = 3
	PENALTY_SG_VALID = 4
	PENALTY_30       = 5
	PENALTY_45       = 6

	// Penatly Reasons
	PENR_UNKNOWN     = 0
	PENR_ADMIN       = 1
	PENR_WRONG_WAY   = 2
	PENR_FALSE_START = 3
	PENR_SPEEDING    = 4
	PENR_STOP_SHORT  = 5
	PENR_STOP_LATE   = 6
)

type Pen struct {
	ReqI uint8 `struct:"uint8"`
	Plid uint8 `struct:"uint8"`

	OldPen uint8 `struct:"uint8"`
	NewPen uint8 `struct:"uint8"`
	Reason uint8 `struct:"uint8"`
	Sp3    uint8 `struct:"uint8"`
}

func (p *Pen) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Pen) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Pen) Type() uint8 {
	return ISP_PLL
}

func NewPen() Packet {
	return &Pen{}
}

func (p *Pen) New() Packet {
	return NewPen()
}
