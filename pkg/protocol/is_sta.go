package protocol

import (
	"encoding/binary"
	"time"

	"github.com/go-restruct/restruct"
)

// IspSta ...
const (
	IspSta = 5
)

// Sta ...
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

// Racing ...
func (p *Sta) Racing() bool {
	return p.RaceInProg == 1
}

// QualifyingDuration ...
func (p *Sta) QualifyingDuration() time.Duration {
	// TODO decode the special rules - should be shared with RST
	return (time.Duration(p.QualMins) * time.Minute)
}

// Laps ...
func (p *Sta) Laps() int32 {
	// TODO decode the special rules - should be shared with RST
	return int32(p.RaceLaps)
}

// UnmarshalInsim ...
func (p *Sta) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Sta) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Sta) Type() uint8 {
	return IspSta
}

// NewSta ...
func NewSta() Packet {
	return &Sta{}
}

// New ...
func (p *Sta) New() Packet {
	return NewSta()
}
