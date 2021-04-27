// +build ignore

package main

import (
	"fmt"
	"github.com/theangryangel/insim-go/pkg/files"
)

func main() {

	p := files.Pth{}
	p.Read("./data/SO1.pth")

	for _, n := range p.Nodes {
		fmt.Printf("C = (%d, %d, %d)\n", n.Centre.X, n.Centre.Y, n.Centre.Z)
	}

	fmt.Printf("NumNodes = %d, Len = %d\n", p.NumNodes, len(p.Nodes))

}
