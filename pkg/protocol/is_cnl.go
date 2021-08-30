package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspCnl ...
const (
	// IspCnl is the packet id
	IspCnl = 19

	CnlReasonDisconnect     = 0
	CnlReasonTimeout        = 1
	CnlReasonLostConnection = 2
	CnlReasonKicked         = 3
	CnlReasonBanned         = 4
	CnlReasonSecurity       = 5
	CnlReasonCPW            = 6
	CnlReasonOOS            = 7
	CnlReasonJOOS           = 8
	CnlReasonHack           = 9
)

// Cnl is a Connection Leave packet
type Cnl struct {
	ReqI uint8 `struct:"uint8"`
	Ucid uint8 `struct:"uint8"`

	Reason uint8 `struct:"uint8"`
	Total  uint8 `struct:"uint8"`
	Spare2 uint8 `struct:"uint8"`
	Spare3 uint8 `struct:"uint8"`
}

// UnmarshalInsim ...
func (p *Cnl) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Cnl) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Cnl) Type() uint8 {
	return IspCnl
}

// NewCnl ...
func NewCnl() Packet {
	return &Cnl{}
}

// New ...
func (p *Cnl) New() Packet {
	return NewCnl()
}
