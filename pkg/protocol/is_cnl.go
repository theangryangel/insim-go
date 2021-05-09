package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_CNL = 19

	LEAVR_DISCO    = 0
	LEAVR_TIMEOUT  = 1
	LEAVR_LOSTCONN = 2
	LEAVR_KICKED   = 3
	LEAVR_BANNED   = 4
	LEAVR_SECURITY = 5
	LEAVR_CPW      = 6
	LEAVR_OOS      = 7
	LEAVR_JOOS     = 8
	LEAVR_HACK     = 9
)

type Cnl struct {
	ReqI uint8 `struct:"uint8"`
	Ucid uint8 `struct:"uint8"`

	Reason uint8 `struct:"uint8"`
	Total  uint8 `struct:"uint8"`
	Spare2 uint8 `struct:"uint8"`
	Spare3 uint8 `struct:"uint8"`
}

func (p *Cnl) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Cnl) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Cnl) Type() uint8 {
	return ISP_NCN
}

func NewCnl() Packet {
	return &Cnl{}
}

func (p *Cnl) New() Packet {
	return NewCnl()
}
