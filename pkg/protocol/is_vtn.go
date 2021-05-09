package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

const (
	ISP_VTN = 16

	VOTE_NONE    = 0
	VOTE_END     = 1
	VOTE_RESTART = 2
	VOTE_QUALIFY = 3
)

type Vtn struct {
	ReqI uint8 `struct:"uint8"`
	Zero uint8 `struct:"uint8"`

	Ucid   uint8 `struct:"uint8"`
	Action uint8 `struct:"uint8"`
	Sp2    uint8 `struct:"uint8"`
	Sp3    uint8 `struct:"uint8"`
}

func (p *Vtn) Voting() bool {
	return p.Action != VOTE_NONE
}

func (p *Vtn) Restart() bool {
	return p.Action == VOTE_RESTART
}

func (p *Vtn) End() bool {
	return p.Action == VOTE_END
}

func (p *Vtn) Qualify() bool {
	return p.Action == VOTE_QUALIFY
}

func (p *Vtn) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

func (p *Vtn) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

func (p *Vtn) Type() uint8 {
	return ISP_VTN
}

func NewVtn() Packet {
	return &Vtn{}
}

func (p *Vtn) New() Packet {
	return NewVtn()
}
