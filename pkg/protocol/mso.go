package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"

	"github.com/theangryangel/insim-go/pkg/strings"
)

const (
	ISP_MSO = 11
)

type Mso struct {
	ReqI uint8 `struct:"uint8"`
	Zero uint8 `struct:"uint8"`

	Ucid      uint8 `struct:"uint8"`
	Plid      uint8 `struct:"uint8"`
	UserType  uint8 `struct:"uint8"`
	TextStart uint8 `struct:"uint8"`
	Msg       string
}

func (p *Mso) UnmarshalInsim(data []byte) (err error) {
	err = restruct.Unpack(data, binary.LittleEndian, p)
	if err != nil {
		return err
	}

	start := p.TextStart + 6
	p.Msg, err = strings.Decode(data[start:])
	//p.Msg = string(data[start:])
	if err != nil {
		return err
	}
	return err
}

func (p *Mso) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Mso) Type() uint8 {
	return ISP_MSO
}

func NewMso() Packet {
	return &Mso{}
}

func (p *Mso) New() Packet {
	return NewMso()
}
