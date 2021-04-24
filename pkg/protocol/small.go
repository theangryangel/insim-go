package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_SMALL = 4

	SMALL_NONE = 0
	SMALL_SSP  = 1
	SMALL_SSG  = 2
	SMALL_VTA  = 3
	SMALL_TMS  = 4
	SMALL_STP  = 5
	SMALL_RTP  = 6
	SMALL_NLI  = 7
	SMALL_ALC  = 8
	SMALL_LCS  = 9
)

type Small struct {
	ReqI uint8 `struct:"uint8"`
	SubT uint8 `struct:"uint8"`

	UVal uint32 `struct:"uint32"`
}

func (p *Small) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Small) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Small) Type() (id uint8) {
	return ISP_SMALL
}

func NewSmall() Packet {
	return &Small{}
}

func (p *Small) New() Packet {
	return NewSmall()
}
