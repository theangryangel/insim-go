// +build ignore

package main

import (
	"flag"
	"fmt"
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

	fmt.Println(p.Svg(
		*imageWidth,
		*imageHeight,
		*resolution,
		*trackcolour,
		*limitcolour,
		*linecolour,
		*startfinishcolour,
		*line,
	))
}
