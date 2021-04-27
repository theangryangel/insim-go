package files

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"io/ioutil"
)

type PthCentre struct {
	X int32 `struct:"int32"`
	Y int32 `struct:"int32"`
	Z int32 `struct:"int32"`
}

type PthDirection struct {
	X float32 `struct:"float32"`
	Y float32 `struct:"float32"`
	Z float32 `struct:"float32"`
}

type PthLimit struct {
	Left  float32 `struct:"float32"`
	Right float32 `struct:"float32"`
}

type PthNode struct {
	Centre    PthCentre
	Direction PthDirection

	OuterLimit PthLimit
	RoadLimit  PthLimit
}

type Pth struct {
	Magic    string `struct:"[6]byte"`
	Version  uint8  `struct:"uint8"`
	Revision uint8  `struct:"uint8"`

	NumNodes   int32 `struct:"int32,sizeof=Nodes"`
	FinishLine int32 `struct:"int32"`

	Nodes []PthNode
}

func (p *Pth) Read(file string) (err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return restruct.Unpack(data, binary.LittleEndian, p)
}
