package protocol

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
)

// IrpHlr ...
const (
	IrpHlr = 252 // Hostlist request
	IrpHos = 253 // Hostlist response
	IrpSel = 254 // Send : To select a host
	IrpErr = 255 // Error

	RelayErrPacket1  = 1 // Invalid packet sent by client (wrong structure / length)
	RelayErrPacket2  = 2 // Invalid packet sent by client (packet was not allowed to be forwarded to host)
	RelayErrHostname = 3 // Wrong hostname given by client
	RelayErrAdmin    = 4 // Wrong admin pass given by client
	RelayErrSpec     = 5 // Wrong spec pass given by client
	RelayErrNoSpec   = 6 // Spectator pass required, but none given
)

// RelayHlr ...
type RelayHlr struct {
	ReqI  uint8 `struct:"uint8"`
	Spare uint8 `struct:"uint8"`
}

// UnmarshalInsim ...
func (p *RelayHlr) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *RelayHlr) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *RelayHlr) Type() uint8 {
	return IrpHlr
}

// NewRelayHlr ...
func NewRelayHlr() Packet {
	return &RelayHlr{}
}

// New ...
func (p *RelayHlr) New() Packet {
	return NewRelayHlr()
}

// RelayHosInfo ...
type RelayHosInfo struct {
	HName    string `struct:"[32]byte"`
	Track    string `struct:"[6]byte"`
	Flags    uint8  `struct:"uint8"`
	NumConns uint8  `struct:"uint8"`
}

// RelayHos ...
type RelayHos struct {
	ReqI     uint8 `struct:"uint8"`
	NumHosts uint8 `struct:"uint8,sizeof=HInfo"`

	HInfo []RelayHosInfo
}

// UnmarshalInsim ...
func (p *RelayHos) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *RelayHos) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *RelayHos) Type() uint8 {
	return IrpHos
}

// NewRelayHos ...
func NewRelayHos() Packet {
	return &RelayHos{}
}

// New ...
func (p *RelayHos) New() Packet {
	return NewRelayHos()
}

// RelaySel ...
type RelaySel struct {
	ReqI uint8 `struct:"uint8"`
	Zero uint8 `struct:"uint8"`

	HName string `struct:"[32]byte"`
	Admin string `struct:"[16]byte"`
	Spec  string `struct:"[16]byte"`
}

// UnmarshalInsim ...
func (p *RelaySel) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *RelaySel) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *RelaySel) Type() uint8 {
	return IrpSel
}

// NewRelaySel ...
func NewRelaySel() Packet {
	return &RelaySel{}
}

// New ...
func (p *RelaySel) New() Packet {
	return NewRelaySel()
}

// RelayErr ...
type RelayErr struct {
	ReqI uint8 `struct:"uint8"`
	Err  uint8 `struct:"uint8"`
}

// ErrMessage ...
func (p *RelayErr) ErrMessage() string {
	switch p.Err {
	case RelayErrPacket1:
		return "Invalid packet sent by client (wrong structure / length)"
	case RelayErrPacket2:
		return "Invalid packet sent by client (packet was not allowed to be forwarded to host)"
	case RelayErrHostname:
		return "Wrong hostname given by client"
	case RelayErrAdmin:
		return "Wrong admin pass given by client"
	case RelayErrSpec:
		return "Wrong spec pass given by client"
	case RelayErrNoSpec:
		return "Spectator pass required, but none given"
	default:
		return "Unknown"
	}
}

// UnmarshalInsim ...
func (p *RelayErr) UnmarshalInsim(data []byte) (err error) {
	return restruct.Unpack(data, binary.LittleEndian, p)
}

// MarshalInsim ...
func (p *RelayErr) MarshalInsim() (data []byte, err error) {
	return restruct.Pack(binary.LittleEndian, p)
}

// Type ...
func (p *RelayErr) Type() uint8 {
	return IrpErr
}

// NewRelayErr ...
func NewRelayErr() Packet {
	return &RelayErr{}
}

// New ...
func (p *RelayErr) New() Packet {
	return NewRelayErr()
}
