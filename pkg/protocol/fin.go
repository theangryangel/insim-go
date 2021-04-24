package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"time"
)

const (
	ISP_FIN = 34
)

type Fin struct {
	ReqI uint8 `struct:"uint8"`

	Plid uint8 `struct:"uint8"`

	TTime uint32 `struct:"uint32"`
	BTime uint32 `struct:"uint32"`

	SpA      uint8 `struct:"uint8"`
	NumStops uint8 `struct:"uint8"`
	Confirm  uint8 `struct:"uint8"`
	SpB      uint8 `struct:"uint8"`

	LapsDone uint16 `struct:"uint16"`
	Flags    uint16 `struct:"uint16"`
}

func (p *Fin) TotalTime() time.Duration {
	return (time.Duration(p.TTime) * time.Millisecond)
}

func (p *Fin) BestTime() time.Duration {
	return (time.Duration(p.BTime) * time.Millisecond)
}

func (p *Fin) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Fin) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Fin) Type() (id uint8) {
	return ISP_FIN
}

func NewFin() Packet {
	return &Fin{}
}

func (p *Fin) New() Packet {
	return NewFin()
}
