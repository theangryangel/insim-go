package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"github.com/theangryangel/insim-go/pkg/strings"
)

// IspIsm ...
const (
	IspIsm = 10
)

// Ism ...
type Ism struct {
	ReqI uint8 `struct:"uint8"`
	Zero uint8 `struct:"uint8"`

	Host uint8 `struct:"uint8"`
	Sp1  uint8 `struct:"uint8"`
	Sp2  uint8 `struct:"uint8"`
	Sp3  uint8 `struct:"uint8"`

	HName string `struct:"[]byte"`
}

// UnmarshalInsim ...
func (p *Ism) UnmarshalInsim(data []byte) (err error) {
	err = restruct.Unpack(data, binary.LittleEndian, p)
	if err != nil {
		return err
	}

	p.HName, err = strings.Decode([]byte(p.HName))
	if err != nil {
		return err
	}
	return nil
}

// MarshalInsim ...
func (p *Ism) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Ism) Type() uint8 {
	return IspIsm
}

// NewIsm ...
func NewIsm() Packet {
	return &Ism{}
}

// New ...
func (p *Ism) New() Packet {
	return NewIsm()
}
