package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_PLA = 28

	PITLANE_EXIT       = 0
	PITLANE_ENTER      = 1
	PITLANE_NO_PURPOSE = 2
	PITLANE_DT         = 3
	PITLANE_SG         = 4
)

type Pla struct {
	ReqI uint8 `struct:"uint8"`
	Plid uint8 `struct:"uint8"`
	Fact uint8 `struct:"uint8"`
	Sp1  uint8 `struct:"uint8"`
	Sp2  uint8 `struct:"uint8"`
	Sp3  uint8 `struct:"uint8"`
}

func (p *Pla) Entering() bool {
	return ((p.Fact >> PITLANE_ENTER) & 1) != 0
}

func (p *Pla) Exiting() bool {
	return ((p.Fact >> PITLANE_EXIT) & 1) != 0
}

func (p *Pla) NoPurpose() bool {
	return ((p.Fact >> PITLANE_NO_PURPOSE) & 1) != 0
}

func (p *Pla) StopGo() bool {
	return ((p.Fact >> PITLANE_SG) & 1) != 0
}

func (p *Pla) DriveThrough() bool {
	return ((p.Fact >> PITLANE_DT) & 1) != 0
}

func (p *Pla) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Pla) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Pla) Type() uint8 {
	return ISP_PLL
}

func NewPla() Packet {
	return &Pla{}
}

func (p *Pla) New() Packet {
	return NewPla()
}
