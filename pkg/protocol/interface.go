package protocol

type Packet interface {
	Unmarshal(data []byte) (err error)
	Marshal() (data []byte, err error)
	Type() (id uint8)
	New() (Packet)
}
