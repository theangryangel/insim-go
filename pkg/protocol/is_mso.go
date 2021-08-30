package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"

	"github.com/theangryangel/insim-go/pkg/strings"
)

// IspMso ...
const (
	IspMso = 11
)

// Mso ...
type Mso struct {
	ReqI uint8 `struct:"uint8"`
	Zero uint8 `struct:"uint8"`

	Ucid      uint8 `struct:"uint8"`
	Plid      uint8 `struct:"uint8"`
	UserType  uint8 `struct:"uint8"`
	TextStart uint8 `struct:"uint8"`
	Msg       string
}

// UnmarshalInsim ...
func (p *Mso) UnmarshalInsim(data []byte) (err error) {
	err = restruct.Unpack(data, binary.LittleEndian, p)
	if err != nil {
		return err
	}

	start := p.TextStart + 6
	p.Msg, err = strings.Decode(data[start:])
	if err != nil {
		return err
	}
	return err
}

// MarshalInsim ...
func (p *Mso) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Mso) Type() uint8 {
	return IspMso
}

// NewMso ...
func NewMso() Packet {
	return &Mso{}
}

// New ...
func (p *Mso) New() Packet {
	return NewMso()
}
