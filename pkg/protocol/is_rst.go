package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"time"
)

const (
	ISP_RST = 17
)

type Rst struct {
	ReqI uint8 `struct:"uint8"`
	Zero uint8 `struct:"uint8"`

	RaceLaps uint8 `struct:"uint8"`
	QualMins uint8 `struct:"uint8"`

	NumP   uint8 `struct:"uint8"`
	Timing uint8 `struct:"uint8"`

	Track   string `struct:"[6]byte"`
	Weather uint8  `struct:"uint8"`
	Wind    uint8  `struct:"uint8"`

	Flags    uint16 `struct:"uint16"`
	NumNodes uint16 `struct:"uint16"`
	Finish   uint16 `struct:"uint16"`
	Split1   uint16 `struct:"uint16"`
	Split2   uint16 `struct:"uint16"`
	Split3   uint16 `struct:"uint16"`
}

func (p *Rst) Qualifying() bool {
	return p.RaceLaps == 0
}

func (p *Rst) Racing() bool {
	return p.QualMins == 0
}

func (p *Rst) QualifyingDuration() time.Duration {
	// TODO decode the special rules - should be shared with RST
	return (time.Duration(p.QualMins) * time.Minute)
}

func (p *Rst) Laps() int32 {
	// TODO decode the special rules - should be shared with RST
	return int32(p.RaceLaps)
}

func (p *Rst) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Rst) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Rst) Type() uint8 {
	return ISP_RST
}

func NewRst() Packet {
	return &Rst{}
}

func (p *Rst) New() Packet {
	return NewRst()
}
