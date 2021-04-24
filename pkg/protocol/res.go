package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"time"
)

const (
	ISP_RES = 35
)

type Res struct {
	ReqI uint8 `struct:"uint8"`

	Plid uint8 `struct:"uint8"`

	UName string `struct:"[24]byte"`
	PName string `struct:"[24]byte"`
	Plate string `struct:"[8]byte"`
	CName string `struct:"[4]byte"`

	TTime uint32 `struct:"uint32"`
	BTime uint32 `struct:"uint32"`

	SpA      uint8 `struct:"uint8"`
	NumStops uint8 `struct:"uint8"`
	Confirm  uint8 `struct:"uint8"`
	SpB      uint8 `struct:"uint8"`

	LapsDone uint16 `struct:"uint16"`
	Flags    uint16 `struct:"uint16"`

	ResultNum uint8 `struct:"uint8"`
	NumRes    uint8 `struct:"uint8"`

	PSeconds uint16 `struct:"uint16"`
}

func (p *Res) TotalTime() time.Duration {
	return (time.Duration(p.TTime) * time.Millisecond)
}

func (p *Res) BestTime() time.Duration {
	return (time.Duration(p.BTime) * time.Millisecond)
}

func (p *Res) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Res) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Res) Type() (id uint8) {
	return ISP_RES
}

func NewRes() Packet {
	return &Res{}
}

func (p *Res) New() Packet {
	return NewRes()
}
