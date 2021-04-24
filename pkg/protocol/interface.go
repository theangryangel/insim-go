package protocol

type Packet interface {
	UnmarshalInsim(data []byte) (err error)
	MarshalInsim() (data []byte, err error)
	Type() (id uint8)
	New() Packet
}
