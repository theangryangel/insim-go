package files

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"io/ioutil"
	"math"
)

var LeftCos = math.Cos(90 * math.Pi / 180)
var LeftSin = math.Sin(90 * math.Pi / 180)
var RightCos = math.Cos(-90 * math.Pi / 180)
var RightSin = math.Sin(-90 * math.Pi / 180)

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

func (node *PthNode) RoadLimits(metres bool) (float64, float64, float64, float64) {
	factor := float64(1)
	if metres {
		factor = 65536
	}

	rlx := (float64(node.Direction.X)*LeftCos-(-float64(node.Direction.Y))*LeftSin)*float64(node.RoadLimit.Left) + (float64(node.Centre.X) / factor)
	rly := (float64(-node.Direction.Y)*LeftCos+float64(node.Direction.X)*LeftSin)*float64(node.RoadLimit.Left) + (float64(-node.Centre.Y) / factor)

	rrx := (float64(node.Direction.X)*RightCos-(-float64(node.Direction.Y))*RightSin)*float64(-node.RoadLimit.Right) + (float64(node.Centre.X) / factor)
	rry := (float64(-node.Direction.Y)*RightCos+float64(node.Direction.X)*RightSin)*float64(-node.RoadLimit.Right) + (float64(-node.Centre.Y) / factor)

	return rlx, rly, rrx, rry
}

func (node *PthNode) OuterLimits(metres bool) (float64, float64, float64, float64) {
	factor := float64(1)
	if metres {
		factor = 65536
	}

	llx := (float64(node.Direction.X)*LeftCos-(-float64(node.Direction.Y))*LeftSin)*float64(node.OuterLimit.Left) + (float64(node.Centre.X) / factor)
	lly := (float64(-node.Direction.Y)*LeftCos+float64(node.Direction.X)*LeftSin)*float64(node.OuterLimit.Left) + (float64(-node.Centre.Y) / factor)

	lrx := (float64(node.Direction.X)*RightCos-(-float64(node.Direction.Y))*RightSin)*float64(-node.OuterLimit.Right) + (float64(node.Centre.X) / factor)
	lry := (float64(-node.Direction.Y)*RightCos+float64(node.Direction.X)*RightSin)*float64(-node.OuterLimit.Right) + (float64(-node.Centre.Y) / factor)

	return llx, lly, lrx, lry
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
