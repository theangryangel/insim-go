package protocol

import (
	"encoding/binary"
	"time"

	"github.com/go-restruct/restruct"
)

const (
	ISP_STA = 5
)

type Sta struct {
	ReqI uint8 `struct:"uint8"`
	Zero uint8 `struct:"uint8"`

	ReplaySpeed float32 `struct:"float32"`
	Flags       uint16  `struct:"uint16"`
	InGameCam   uint8   `struct:"uint8"`
	ViewPlid    uint8   `struct:"uint8"`
	NumP        uint8   `struct:"uint8"`
	NumConns    uint8   `struct:"uint8"`
	NumFinished uint8   `struct:"uint8"`
	RaceInProg  uint8   `struct:"uint8"`
	QualMins    uint8   `struct:"uint8"`
	RaceLaps    uint8   `struct:"uint8"`
	Spare2      uint8   `struct:"uint8"`
	Spare3      uint8   `struct:"uint8"`
	Track       string  `struct:"[6]byte"`
	Weather     uint8   `struct:"uint8"`
	Wind        uint8   `struct:"uint8"`
}

func (p *Sta) Racing() bool {
	return p.RaceInProg == 1
}

func (p *Sta) QualifyingDuration() time.Duration {
	// TODO decode the special rules
	return (time.Duration(p.QualMins) * time.Minute)
}

func (p *Sta) Unmarshal(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Sta) Marshal() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Sta) Type() uint8 {
	return ISP_STA
}

func NewSta() Packet {
	return &Sta{}
}

func (p *Sta) New() Packet {
	return NewSta()
}
