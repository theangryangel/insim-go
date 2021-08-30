package protocol

import (
	"encoding/binary"

	"github.com/go-restruct/restruct"
	"github.com/theangryangel/insim-go/pkg/strings"
)

// IspNcn ...
const (
	IspNcn = 18
)

// Ncn ...
type Ncn struct {
	ReqI uint8 `struct:"uint8"`
	Ucid uint8 `struct:"uint8"`

	UName string `struct:"[24]byte"`
	PName string `struct:"[24]byte"`
	Admin uint8  `struct:"uint8"`
	Total uint8  `struct:"uint8"`
	Flags uint8  `struct:"uint8"`
	Spare uint8  `struct:"uint8"`
}

// IsAdmin ...
func (p *Ncn) IsAdmin() bool {
	return p.Admin == 1
}

// IsRemote ...
func (p *Ncn) IsRemote() bool {
	return false // TODO
}

// UnmarshalInsim ...
func (p *Ncn) UnmarshalInsim(data []byte) (err error) {
	err = restruct.Unpack(data, binary.LittleEndian, p)
	if err != nil {
		return err
	}

	p.PName, err = strings.Decode([]byte(p.PName))
	if err != nil {
		return err
	}

	p.UName, err = strings.Decode([]byte(p.UName))
	if err != nil {
		return err
	}

	return nil
}

// MarshalInsim ...
func (p *Ncn) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Ncn) Type() uint8 {
	return IspNcn
}

// NewNcn ...
func NewNcn() Packet {
	return &Ncn{}
}

// New ...
func (p *Ncn) New() Packet {
	return NewNcn()
}
