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

	minX, minY, maxX, maxY := p.OuterBounds(true)

	for i, j := 0, len(p.Nodes)-1; i < len(p.Nodes) && j > 0; i, j = i+1, j-1 {

		if *resolution > 0 && i%*resolution != 0 {
			continue
		}

		node := &p.Nodes[i]

		rcx, rcy := node.RoadCentre(true)

		roadCX = append(roadCX, rcx)
		roadCY = append(roadCY, rcy)

		// calc road
		rlx, rly, rrx, rry := node.RoadLimits(true)

		roadLX = append(roadLX, rlx)
		roadLY = append(roadLY, rly)

		roadRX = append(roadRX, rrx)
		roadRY = append(roadRY, rry)

		// calc limit
		llx, lly, lrx, lry := node.OuterLimits(true)

		outerLX = append(outerLX, llx)
		outerLY = append(outerLY, lly)

		outerRX = append(outerRX, lrx)
		outerRY = append(outerRY, lry)
	}

	lastnode := p.Nodes[0]

	rcx, rcy := lastnode.RoadCentre(true)

	roadCX = append(roadCX, rcx)
	roadCY = append(roadCY, rcy)

	rlx, rly, rrx, rry := lastnode.RoadLimits(true)

	roadLX = append(roadLX, rlx)
	roadLY = append(roadLY, rly)

	roadRX = append(roadRX, rrx)
	roadRY = append(roadRY, rry)

	// calc limit
	llx, lly, lrx, lry := lastnode.OuterLimits(true)

	outerLX = append(outerLX, llx)
	outerLY = append(outerLY, lly)

	outerRX = append(outerRX, lrx)
	outerRY = append(outerRY, lry)

	s := svg.New(os.Stdout)

	disX := maxX - minX
	disY := maxY - minY

	// we use viewbox to scale so we dont have to scale the coordinates by hand
	// this means our text will be wonky sized. boo.
	viewBox := fmt.Sprintf(
		"viewBox=\"%.2f %.2f %.2f %.2f\"",
		minX-(0.01*disX),
		minY-(0.01*disY),
		disX+(0.02*disX),
		disY+(0.02*disY),
	)

	s.Start(*imageWidth, *imageHeight, viewBox)

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
