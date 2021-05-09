package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"

	"github.com/theangryangel/insim-go/pkg/strings"
)

const (
	ISP_NPL = 21
)

type Npl struct {
	ReqI uint8 `struct:"uint8"`

	Plid  uint8    `struct:"uint8"`
	Ucid  uint8    `struct:"uint8"`
	PType uint8    `struct:"uint8"`
	Flags uint16   `struct:"uint16"`
	PName string   `struct:"[24]byte"`
	Plate string   `struct:"[8]byte"`
	CName string   `struct:"[4]byte"`
	SName string   `struct:"[16]byte"`
	Tyres [4]uint8 `struct:"[4]uint8"`
	HMass uint8    `struct:"uint8"`
	HTres uint8    `struct:"uint8"`
	Model uint8    `struct:"uint8"`
	Pass  uint8    `struct:"uint8"`

	Spare int32 `struct:"int32"`

	SetF   uint8 `struct:"uint8"`
	NumP   uint8 `struct:"uint8"`
	Spare2 uint8 `struct:"uint8"`
	Spare3 uint8 `struct:"uint8"`
}

func (p *Npl) UnmarshalInsim(data []byte) (err error) {
	err = restruct.Unpack(data, binary.LittleEndian, p)
	if err != nil {
		return err
	}

	p.PName, err = strings.Decode([]byte(p.PName))
	if err != nil {
		return err
	}

	p.Plate, err = strings.Decode([]byte(p.Plate))
	if err != nil {
		return err
	}

	return nil
}

func (p *Npl) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Npl) Type() uint8 {
	return ISP_NCN
}

func NewNpl() Packet {
	return &Npl{}
}

func (p *Npl) New() Packet {
	return NewNpl()
}
