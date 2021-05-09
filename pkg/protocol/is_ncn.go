package protocol

import (
	"encoding/binary"

	"github.com/go-restruct/restruct"
	"github.com/theangryangel/insim-go/pkg/strings"
)

const (
	ISP_NCN = 18
)

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

func (p *Ncn) IsAdmin() bool {
	return p.Admin == 1
}

func (p *Ncn) IsRemote() bool {
	return false // TODO
}

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

func (p *Ncn) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Ncn) Type() uint8 {
	return ISP_NCN
}

func NewNcn() Packet {
	return &Ncn{}
}

func (p *Ncn) New() Packet {
	return NewNcn()
}
