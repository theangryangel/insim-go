package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspFlg ...
const (
	IspFlg = 32

	FlagYellow = 1
	FlagBlue   = 2
)

// Flg ...
type Flg struct {
	ReqI uint8 `struct:"uint8"`
	Plid uint8 `struct:"uint8"`

	OffOn     uint8 `struct:"uint8"`
	Flag      uint8 `struct:"uint8"`
	CarBehind uint8 `struct:"uint8"`
	Sp3       uint8 `struct:"uint8"`
}

// Changed ...
func (p *Flg) Changed() (int, bool) {
	return int(p.Flag), (p.OffOn > 0)
}

// UnmarshalInsim ...
func (p *Flg) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Flg) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Flg) Type() (id uint8) {
	return IspFlg
}

// NewFlg ...
func NewFlg() Packet {
	return &Flg{}
}

// New ...
func (p *Flg) New() Packet {
	return NewFlg()
}
