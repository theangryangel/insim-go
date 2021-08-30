package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspIsi ...
const (
	IspIsi = 1
)

// Init ...
type Init struct {
	ReqI     uint8  `struct:"uint8"`
	Zero     uint8  `struct:"uint8"`
	UDPPort  uint   `struct:"uint16"`
	Flags    uint16 `struct:"uint16"`
	InSimVer uint8  `struct:"uint8"`
	Prefix   byte   `struct:"uint8"`
	Interval uint16 `struct:"uint16"`
	Admin    string `struct:"[16]byte"`
	IName    string `struct:"[16]byte"`
}

// UnmarshalInsim ...
func (p *Init) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Init) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Init) Type() (id uint8) {
	return IspIsi
}

// New ...
func (p *Init) New() Packet {
	return NewInit()
}

// NewInit ...
func NewInit() Packet {
	flags := uint16(0)
	flags = flags | 32

	return &Init{
		InSimVer: 8,
		IName:    "insim.go",
		Interval: 1000,
		Flags:    flags,
	}
}
