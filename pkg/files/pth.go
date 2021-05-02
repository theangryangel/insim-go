package files

import (
	"encoding/binary"
	"github.com/go-restruct/restruct"
	"github.com/theangryangel/insim-go/pkg/geometry"
	"io/ioutil"
	"math"
)

var LeftCos = math.Cos(90 * math.Pi / 180)
var LeftSin = math.Sin(90 * math.Pi / 180)
var RightCos = math.Cos(-90 * math.Pi / 180)
var RightSin = math.Sin(-90 * math.Pi / 180)

type PthLimit struct {
	Left  float32 `struct:"float32"`
	Right float32 `struct:"float32"`
}

type PthNode struct {
	Centre    geometry.FixedPoint
	Direction geometry.FloatingPoint

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
	Magic    string `struct:"[6]byte" json:"-"`
	Version  uint8  `struct:"uint8" json:"-"`
	Revision uint8  `struct:"uint8" json:"-"`

	NumNodes   int32 `struct:"int32,sizeof=Nodes"`
	FinishLine int32 `struct:"int32"`

	Nodes []PthNode
}

func (p *Pth) Read(file string) (err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = restruct.Unpack(data, binary.LittleEndian, p)
	if err != nil {
		return err
	}

	return nil
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

type PthFit struct {
	// TODO reuse Pth?

	OuterX []float64
	OuterY []float64

	RoadX []float64
	RoadY []float64

	RoadCX []float64
	RoadCY []float64

	FinishX []float64
	FinishY []float64

	ScaleX float64
	ScaleY float64

	TranslateX float64
	TranslateY float64
}

func (p *Pth) FitTo(imageWidth float64, imageHeight float64, resolution int) PthFit {
	// TODO refactor a bit
	var roadCX []float64
	var roadCY []float64

	var roadLX []float64
	var roadLY []float64

	var roadRX []float64
	var roadRY []float64

	var outerLX []float64
	var outerLY []float64

	var outerRX []float64
	var outerRY []float64

	scalex, scaley, translatex, translatey := p.GenerateScaleTransform(imageWidth, imageHeight)

	for i := 0; i < len(p.Nodes); i++ {

		if resolution > 0 && i%resolution != 0 {
			continue
		}

		node := &p.Nodes[i]

		rcx, rcy := node.RoadCentre(true)

		roadCX = append(roadCX, rcx*scalex+translatex)
		roadCY = append(roadCY, rcy*scaley+translatey)

		// calc road
		rlx, rly, rrx, rry := node.RoadLimits(true)

		roadLX = append(roadLX, rlx*scalex+translatex)
		roadLY = append(roadLY, rly*scaley+translatey)

		roadRX = append(roadRX, rrx*scalex+translatex)
		roadRY = append(roadRY, rry*scaley+translatey)

		// calc limit
		llx, lly, lrx, lry := node.OuterLimits(true)

		outerLX = append(outerLX, llx*scalex+translatex)
		outerLY = append(outerLY, lly*scaley+translatey)

		outerRX = append(outerRX, lrx*scalex+translatex)
		outerRY = append(outerRY, lry*scaley+translatey)
	}

	// copy the first node to close the loop

	roadCX = append(roadCX, roadCX[0])
	roadCY = append(roadCY, roadCY[0])

	roadLX = append(roadLX, roadLX[0])
	roadLY = append(roadLY, roadLY[0])

	roadRX = append(roadRX, roadRX[0])
	roadRY = append(roadRY, roadRY[0])

	outerLX = append(outerLX, outerLX[0])
	outerLY = append(outerLY, outerLY[0])

	outerRX = append(outerRX, outerRX[0])
	outerRY = append(outerRY, outerRY[0])

	flx, fly, frx, fry := p.Nodes[p.FinishLine].RoadLimits(true)

	flx = flx*scalex + translatex
	fly = fly*scaley + translatey
	frx = frx*scalex + translatex
	fry = fry*scaley + translatey

	res := PthFit{
		OuterX: append(outerLX, outerRX...),
		OuterY: append(outerLY, outerRY...),

		RoadX: append(roadLX, roadRX...),
		RoadY: append(roadLY, roadRY...),

		RoadCX: roadCX,
		RoadCY: roadCY,

		FinishX: []float64{flx, frx},
		FinishY: []float64{fly, fry},

		ScaleX: scalex,
		ScaleY: scaley,

		TranslateX: translatex,
		TranslateY: translatey,
	}

	return res
}
