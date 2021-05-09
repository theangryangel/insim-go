package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	IRP_HLR = 252 // Hostlist request
	IRP_HOS = 253 // Hostlist response
	IRP_SEL = 254 // Send : To select a host
	IRP_ERR = 255 // Error

	IR_ERR_PACKET1  = 1 // Invalid packet sent by client (wrong structure / length)
	IR_ERR_PACKET2  = 2 // Invalid packet sent by client (packet was not allowed to be forwarded to host)
	IR_ERR_HOSTNAME = 3 // Wrong hostname given by client
	IR_ERR_ADMIN    = 4 // Wrong admin pass given by client
	IR_ERR_SPEC     = 5 // Wrong spec pass given by client
	IR_ERR_NOSPEC   = 6 // Spectator pass required, but none given
)

type IrpHlr struct {
	ReqI  uint8 `struct:"uint8"`
	Spare uint8 `struct:"uint8"`
}

func (p *IrpHlr) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *IrpHlr) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *IrpHlr) Type() uint8 {
	return IRP_HLR
}

func NewIrpHlr() Packet {
	return &IrpHlr{}
}

func (p *IrpHlr) New() Packet {
	return NewIrpHlr()
}

type IrpHosInfo struct {
	HName    string `struct:"[32]byte"`
	Track    string `struct:"[6]byte"`
	Flags    uint8  `struct:"uint8"`
	NumConns uint8  `struct:"uint8"`
}

type IrpHos struct {
	ReqI     uint8 `struct:"uint8"`
	NumHosts uint8 `struct:"uint8,sizeof=HInfo"`

	HInfo []IrpHosInfo
}

func (p *IrpHos) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *IrpHos) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *IrpHos) Type() uint8 {
	return IRP_HOS
}

func NewIrpHos() Packet {
	return &IrpHos{}
}

func (p *IrpHos) New() Packet {
	return NewIrpHos()
}

type IrpSel struct {
	ReqI uint8 `struct:"uint8"`
	Zero uint8 `struct:"uint8"`

	HName string `struct:"[32]byte"`
	Admin string `struct:"[16]byte"`
	Spec  string `struct:"[16]byte"`
}

func (p *IrpSel) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *IrpSel) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *IrpSel) Type() uint8 {
	return IRP_SEL
}

func NewIrpSel() Packet {
	return &IrpSel{}
}

func (p *IrpSel) New() Packet {
	return NewIrpSel()
}

type IrpErr struct {
	ReqI uint8 `struct:"uint8"`
	Err  uint8 `struct:"uint8"`
}

func (p *IrpErr) ErrMessage() string {
	switch p.Err {
	case IR_ERR_PACKET1:
		return "Invalid packet sent by client (wrong structure / length)"
	case IR_ERR_PACKET2:
		return "Invalid packet sent by client (packet was not allowed to be forwarded to host)"
	case IR_ERR_HOSTNAME:
		return "Wrong hostname given by client"
	case IR_ERR_ADMIN:
		return "Wrong admin pass given by client"
	case IR_ERR_SPEC:
		return "Wrong spec pass given by client"
	case IR_ERR_NOSPEC:
		return "Spectator pass required, but none given"
	default:
		return "Unknown"
	}
}

func (p *IrpErr) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *IrpErr) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *IrpErr) Type() uint8 {
	return IRP_ERR
}

func NewIrpErr() Packet {
	return &IrpErr{}
}

func (p *IrpErr) New() Packet {
	return NewIrpErr()
}
