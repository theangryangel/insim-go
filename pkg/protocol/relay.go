package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	IRP_SEL	= 254	// Send : To select a host
)

type IrpSel struct {
	ReqI uint8 `struct:"uint8"`
	Zero uint8 `struct:"uint8"`

	HName string `struct:"[32]byte"`
	Admin string `struct:"[16]byte"`
	Spec string `struct:"[16]byte"`
}

func (p *IrpSel) Unmarshal(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *IrpSel) Marshal() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *IrpSel) Type() (uint8) {
	return IRP_SEL
}

func NewIrpSel() (Packet) {
	return &IrpSel{
	}
}

func (p *IrpSel) New() (Packet) {
	return NewIrpSel()
}
