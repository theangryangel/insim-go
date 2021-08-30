package files

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"github.com/theangryangel/insim-go/pkg/geometry"
	"io/ioutil"
	"math"
)

var leftCos = math.Cos(90 * math.Pi / 180)
var leftSin = math.Sin(90 * math.Pi / 180)
var rightCos = math.Cos(-90 * math.Pi / 180)
var rightSin = math.Sin(-90 * math.Pi / 180)

// PthLimit describes the left and right limits of the track, from the "centre" point
// Left and Right is dependent on the direction of track around the circuit.
type PthLimit struct {
	Left  float32 `struct:"float32"`
	Right float32 `struct:"float32"`
}

// PthNode describes a point on the track: the centre point, the direction and the limits of the track and road
// These values are the raw values and should not be trusted. Use RoadCentre, RoadLimits, etc. as these will have
// the correct rotations and offsets applied.
type PthNode struct {
	Centre    geometry.FixedPoint
	Direction geometry.FloatingPoint

	OuterLimit PthLimit
	RoadLimit  PthLimit
}

// RoadCentre will return the X,Y coordinates of the track centre at this node
func (node *PthNode) RoadCentre(metres bool) (float64, float64) {
	factor := float64(1)
	if metres {
		factor = 65536
	}

	return (float64(node.Centre.X) / factor), (float64(-node.Centre.Y) / factor)
}

// RoadLimits will return X, Y coordinates for the track limits, at this node,
// with the correct rotations, etc. applied.
func (node *PthNode) RoadLimits(metres bool) (float64, float64, float64, float64) {
	factor := float64(1)
	if metres {
		factor = 65536
	}

	rlx := (float64(node.Direction.X)*leftCos-(-float64(node.Direction.Y))*leftSin)*float64(node.RoadLimit.Left) + (float64(node.Centre.X) / factor)
	rly := (float64(-node.Direction.Y)*leftCos+float64(node.Direction.X)*leftSin)*float64(node.RoadLimit.Left) + (float64(-node.Centre.Y) / factor)

	rrx := (float64(node.Direction.X)*rightCos-(-float64(node.Direction.Y))*rightSin)*float64(-node.RoadLimit.Right) + (float64(node.Centre.X) / factor)
	rry := (float64(-node.Direction.Y)*rightCos+float64(node.Direction.X)*rightSin)*float64(-node.RoadLimit.Right) + (float64(-node.Centre.Y) / factor)

	return rlx, rly, rrx, rry
}

// OuterLimits will return X, Y coordinates for the outer limits, at this node,
// with the correct rotations, etc. applied.
func (node *PthNode) OuterLimits(metres bool) (float64, float64, float64, float64) {
	factor := float64(1)
	if metres {
		factor = 65536
	}

	llx := (float64(node.Direction.X)*leftCos-(-float64(node.Direction.Y))*leftSin)*float64(node.OuterLimit.Left) + (float64(node.Centre.X) / factor)
	lly := (float64(-node.Direction.Y)*leftCos+float64(node.Direction.X)*leftSin)*float64(node.OuterLimit.Left) + (float64(-node.Centre.Y) / factor)

	lrx := (float64(node.Direction.X)*rightCos-(-float64(node.Direction.Y))*rightSin)*float64(-node.OuterLimit.Right) + (float64(node.Centre.X) / factor)
	lry := (float64(-node.Direction.Y)*rightCos+float64(node.Direction.X)*rightSin)*float64(-node.OuterLimit.Right) + (float64(-node.Centre.Y) / factor)

	return llx, lly, lrx, lry
}

// Pth is a PTH file
type Pth struct {
	Magic    string `struct:"[6]byte" json:"-"`
	Version  uint8  `struct:"uint8" json:"-"`
	Revision uint8  `struct:"uint8" json:"-"`

	NumNodes   int32 `struct:"int32,sizeof=Nodes"`
	FinishLine int32 `struct:"int32"`

	Nodes []PthNode
}

// NewPth will create a new Pth from a file
func NewPth(file string) (*Pth, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	p := &Pth{}

	err = restruct.Unpack(data, binary.LittleEndian, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// OuterBounds will return X, Y coordinates for the outer bounds, at this node,
// with the correct rotations, etc. applied.
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

// ScaleTo will return values to help fit a Pth to a specific size (scalex, scalex and translatex, translatey)
func (p *Pth) ScaleTo(imageWidth float64, imageHeight float64) (float64, float64, float64, float64) {
	minX, minY, maxX, maxY := p.OuterBounds(true)

	disX := maxX - minX
	disY := maxY - minY

	/*
		Let vb-x, vb-y, vb-width, vb-height be the min-x, min-y, width and height values of the viewBox attribute respectively.
	*/

	vbx, vby, vbh, vbw := minX,
		minY,
		disY,
		disX

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
