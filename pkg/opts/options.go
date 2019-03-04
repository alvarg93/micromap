package opts

import (
	"fmt"
	"strings"
)

var help = `

Micromap generates a dependency graph from code annotations.

Usage: micromap [-x=val][-r][-d=val][-i=val][-f=val]
  -h, --help 				- this message
	-x, --regex 			- regex pattern of files to scan
	-r, --recursive 	- scan recursive
	-d, --dot 				- output file for a dot graph
	-i, --img 				- output file for visual representation of the graph
	-f, --format 			- format of the image file
`

//Options represent command line arguments for the program
type Options struct {
	Regex     string
	Recursive bool
	DotFile   string
	ImgFile   string
	ImgFormat string
}

//ParseArgs parses command line arguments into Options.
//It returns Options and a boolean flag indicating if
//help menu was requested.
func ParseArgs(args []string) (opts Options, h bool) {
	o := Options{
		Regex:     ".+",
		DotFile:   "micromap.dot",
		ImgFile:   "micromap.png",
		ImgFormat: "png",
		Recursive: false,
	}

	for _, arg := range args {
		optVal := strings.Split(arg, "=")
		switch optVal[0] {
		case "-x":
		case "--regex":
			o.Regex = optVal[1]
		case "-r":
		case "--recursive":
			o.Recursive = true
		case "-d":
		case "--dot":
			o.DotFile = optVal[1]
		case "-i":
		case "--img":
			o.ImgFile = optVal[1]
		case "-f":
		case "--format":
			o.ImgFormat = optVal[1]
		case "-h":
		case "--help":
			fmt.Print(help)
			return o, true
		default:
		}
	}
	return o, false
}
