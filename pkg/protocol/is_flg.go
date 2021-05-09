package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_FLG = 32

	FLG_YELLOW = 1
	FLG_BLUE   = 2
)

type Flg struct {
	ReqI uint8 `struct:"uint8"`
	Plid uint8 `struct:"uint8"`

	OffOn     uint8 `struct:"uint8"`
	Flag      uint8 `struct:"uint8"`
	CarBehind uint8 `struct:"uint8"`
	Sp3       uint8 `struct:"uint8"`
}

func (p *Flg) Changed() (int, bool) {
	return int(p.Flag), (p.OffOn > 0)
}

func (p *Flg) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Flg) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Flg) Type() (id uint8) {
	return ISP_FLG
}

func NewFlg() Packet {
	return &Flg{}
}

func (p *Flg) New() Packet {
	return NewFlg()
}
