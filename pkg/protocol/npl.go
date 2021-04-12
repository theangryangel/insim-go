package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_NPL = 21
)

type Npl struct {
	ReqI uint8 `struct:"uint8"`
	
	Plid uint8 `struct:"uint8"`
	Ucid uint8 `struct:"uint8"`

	Flags uint16 `struct:"uint16"`

	PName string `struct:"[24]byte"`
	Plate string `struct:"[8]byte"`

	CName string `struct:"[4]byte"`
	SName string `struct:"[16]byte"`
	Tyres [4]uint8 `struct:"[4]uint8"`
	HMass uint8 `struct:"uint8"`
	HTres uint8 `struct:"uint8"`
	Model uint8 `struct:"uint8"`
	Pass uint8 `struct:"uint8"`

	Spare int32 `struct:"int32"`

	SetF uint8 `struct:"uint8"`
	NumP uint8 `struct:"uint8"`
	Spare2 uint8 `struct:"uint8"`
	Spare3 uint8 `struct:"uint8"`
}

func (p *Npl) Unmarshal(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Npl) Marshal() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Npl) Type() (uint8) {
	return ISP_NCN
}

func NewNpl() (Packet) {
	return &Npl{
	}
}

func (p *Npl) New() (Packet) {
	return NewNpl()
}
