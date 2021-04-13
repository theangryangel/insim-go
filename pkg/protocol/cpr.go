package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_CPR = 20
)

type Cpr struct {
	ReqI uint8 `struct:"uint8"`
	Ucid uint8 `struct:"uint8"`

	PName string `struct:"[24]byte"`
	Plate string `struct:"[8]byte"`
}

func (p *Cpr) Unmarshal(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Cpr) Marshal() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Cpr) Type() uint8 {
	return ISP_NCN
}

func NewCpr() Packet {
	return &Cpr{}
}

func (p *Cpr) New() Packet {
	return NewCpr()
}
