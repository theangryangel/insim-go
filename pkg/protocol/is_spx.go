package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"time"
)

const (
	ISP_SPX = 25
)

type Spx struct {
	ReqI uint8 `struct:"uint8"`

	Plid uint8 `struct:"uint8"`

	STime uint32 `struct:"uint32"`
	ETime uint32 `struct:"uint32"`

	Split    uint8 `struct:"uint8"`
	Penalty  uint8 `struct:"uint8"`
	NumStops uint8 `struct:"uint8"`
	Fuel200  uint8 `struct:"uint8"`
}

func (p *Spx) Fuel() uint8 {
	// TODO magic decoding
	return p.Fuel200
}

func (p *Spx) SplitTime() time.Duration {
	return (time.Duration(p.STime) * time.Millisecond)
}

func (p *Spx) ElapsedTime() time.Duration {
	return (time.Duration(p.ETime) * time.Millisecond)
}

func (p *Spx) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Spx) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Spx) Type() (id uint8) {
	return ISP_SPX
}

func NewSpx() Packet {
	return &Spx{}
}

func (p *Spx) New() Packet {
	return NewSpx()
}
