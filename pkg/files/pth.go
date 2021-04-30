package files

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"github.com/theangryangel/insim.go/pkg/geometry"
	"io/ioutil"
	"math"
)

var LeftCos = math.Cos(90 * math.Pi / 180)
var LeftSin = math.Sin(90 * math.Pi / 180)
var RightCos = math.Cos(-90 * math.Pi / 180)
var RightSin = math.Sin(-90 * math.Pi / 180)

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
	Centre    geometry.FixedPoint
	Direction PthDirection

	OuterLimit PthLimit
	RoadLimit  PthLimit
}

func (node *PthNode) RoadCentre(metres bool) (float64, float64) {
	factor := float64(1)
	if metres {
		factor = 65536
	}

	return (float64(node.Centre.X) / factor), (float64(-node.Centre.Y) / factor)
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

func (p *Pth) OuterBounds(metres bool) (float64, float64, float64, float64) {

	minx := float64(0.0)
	miny := float64(0.0)
	maxx := float64(0.0)
	maxy := float64(0.0)

	for _, node := range p.Nodes {
		llx, lly, lrx, lry := node.OuterLimits(metres)

		minx = math.Min(minx, llx)
		minx = math.Min(minx, lrx)

		maxx = math.Max(maxx, llx)
		maxx = math.Max(maxx, lrx)

		miny = math.Min(miny, lly)
		miny = math.Min(miny, lry)

		maxy = math.Max(maxy, lly)
		maxy = math.Max(maxy, lry)
	}

	return minx, miny, maxx, maxy

}

func (p *Pth) GenerateScaleTransform(imageWidth float64, imageHeight float64) (float64, float64, float64, float64) {
	// TODO this should probably be refactored out of Pth, it would be useful in other plces.
	// But we should introduce some general geometry types and funcs

	minX, minY, maxX, maxY := p.OuterBounds(true)

	disX := maxX - minX
	disY := maxY - minY

	/*
		Let vb-x, vb-y, vb-width, vb-height be the min-x, min-y, width and height values of the viewBox attribute respectively.
	*/

	vbx, vby, vbh, vbw := minX-(0.01*disX),
		minY-(0.01*disY),
		disY+(0.02*disY),
		disX+(0.02*disX)

	/*
		Let e-x, e-y, e-width, e-height be the position and size of the element respectively.
	*/
	ex, ey, ew, eh := 0.0, 0.0, imageWidth, imageHeight

	/*
		Initialize scale-x to e-width/vb-width.
		Initialize scale-y to e-height/vb-height.
		Set the larger of scale-x and scale-y to the smaller.
	*/

	scalex := ew / vbw
	scaley := eh / vbh
	if scalex > scaley {
		scalex = scaley
	} else {
		scaley = scalex
	}

	/*
		Initialize translate-x to e-x - (vb-x * scale-x).
		Initialize translate-y to e-y - (vb-y * scale-y)
		If align contains 'xMid', add (e-width - vb-width * scale-x) / 2 to translate-x.
		If align contains 'yMid', add (e-height - vb-height * scale-y) / 2 to translate-y.
	*/

	translatex := (ex - (vbx * scalex)) + ((ew - vbw*scalex) / 2)
	translatey := (ey - (vby * scaley)) + ((eh - vbh*scaley) / 2)

	return scalex, scaley, translatex, translatey
}
