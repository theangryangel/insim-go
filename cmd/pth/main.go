// +build ignore

package main

import (
	"flag"
	"fmt"
	svg "github.com/ajstarks/svgo/float"
	"github.com/theangryangel/insim-go/pkg/files"
	"os"
)

func main() {
	file := flag.String("map", "BL1", "Path file/track")
	resolution := flag.Int("resolution", 1.0, "Resolution. 1=full quality, 0=auto, 2=1/2, 3=1/3, etc.")
	trackcolour := flag.String("track", "#1F2937", "Track colour in hex.")
	limitcolour := flag.String("limit", "#059669", "Limit colour in hex.")
	linecolour := flag.String("line", "#F9FAFB", "Racing line colour")
	line := flag.Bool("showline", true, "Show racing line")
	startfinishcolour := flag.String("startfinish", "#ffffff", "Start/Finish colour in hex.")
	imageHeight := flag.Float64("height", 1024, "Image width")
	imageWidth := flag.Float64("width", 1024, "Image Height")

	flag.Parse()

	p := files.Pth{}

	if _, err := os.Stat(*file); os.IsNotExist(err) {
		panic(err)
	}

	p.Read(*file)

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

	scalex, scaley, translatex, translatey := p.GenerateScaleTransform(*imageWidth, *imageHeight)

	for i := 0; i < len(p.Nodes); i++ {

		if *resolution > 0 && i%*resolution != 0 {
			continue
		}

		node := &p.Nodes[i]

		rcx, rcy := node.RoadCentre(true)

		roadCX = append(roadCX, rcx * scalex * translatex)
		roadCY = append(roadCY, rcy * scaley * translatey)

		// calc road
		rlx, rly, rrx, rry := node.RoadLimits(true)

		roadLX = append(roadLX, rlx * scalex + translatex)
		roadLY = append(roadLY, rly * scaley + translatey)

		roadRX = append(roadRX, rrx * scalex + translatex)
		roadRY = append(roadRY, rry * scaley + translatey)

		// calc limit
		llx, lly, lrx, lry := node.OuterLimits(true)

		outerLX = append(outerLX, llx * scalex + translatex)
		outerLY = append(outerLY, lly * scaley + translatey)

		outerRX = append(outerRX, lrx * scalex + translatex)
		outerRY = append(outerRY, lry * scaley + translatey)
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

	s := svg.New(os.Stdout)

	s.Start(*imageWidth, *imageHeight, "style=\"border: 1px solid red\"")

	s.Polygon(
		append(outerLX, outerRX...),
		append(outerLY, outerRY...),
		fmt.Sprintf("stroke: %s; stroke-width:2px; fill: %s; fill-rule: evenodd", *limitcolour, *limitcolour),
	)
	s.Polygon(
		append(roadLX, roadRX...),
		append(roadLY, roadRY...),
		fmt.Sprintf("stroke: %s; stroke-width:2px; fill: %s; fill-rule: evenodd", *trackcolour, *trackcolour),
	)

	flx, fly, frx, fry := p.Nodes[p.FinishLine].RoadLimits(true)

	flx = flx * scalex + translatex
	fly = fly * scaley + translatey
	frx = frx * scalex + translatex
	fry = fry * scalex + translatey

	s.Line(
		flx, fly, frx, fry,
		fmt.Sprintf("stroke: %s; stroke-width: 2px;", *startfinishcolour),
	)

	flagX := flx
	flagY := fly

	if flx > frx {
		flagX = frx
		flagY = fry
	}

	s.Text(flagX, flagY, "🏁", " font-size: 25px;", "transform=\"translate(-25, 0)\"")

	if *line {
		s.Polyline(
			roadCX,
			roadCY,
			fmt.Sprintf("stroke: %s; stroke-width: 2px; fill: none", *linecolour),
		)
	}

	s.End()
}
