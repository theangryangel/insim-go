package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IspSmall ...
const (
	IspSmall = 4

	SmallNone = 0
	SmallSSP  = 1
	SmallSSG  = 2
	SmallVTA  = 3
	SmallTMS  = 4
	SmallSTP  = 5
	SmallRTP  = 6
	SmallNLI  = 7
	SmallALC  = 8
	SmallLCS  = 9
)

// Small ...
type Small struct {
	ReqI uint8 `struct:"uint8"`
	SubT uint8 `struct:"uint8"`

	UVal uint32 `struct:"uint32"`
}

// UnmarshalInsim ...
func (p *Small) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *Small) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *Small) Type() (id uint8) {
	return IspSmall
}

// NewSmall ...
func NewSmall() Packet {
	return &Small{}
}

// New ...
func (p *Small) New() Packet {
	return NewSmall()
}
