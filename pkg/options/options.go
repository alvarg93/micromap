package options

import (
	"fmt"
	"strings"
)

var micromapHelp = `
Micromap generates a dependency graph from a specification file in yaml.
Dependencies: dot (graphviz package)
Usage: mmap -y=yaml_file [-d=dot_file] [-f=image_format] [-i=image_file] [-s=stylesheet_file]

	-d, --dot         - output file for a dot graph (default: micromap.dot)
	-f, --format      - format of the image file (default: png)
	-h, --help        - this message
	-i, --img         - output file for a visual representation of the graph (default: micromap.png)
	-s, --style       - stylesheet file (default: micromap.style)
	-y, --yaml        - input file (default: micromap.yml)
`

//Options represent command line arguments for the program
type Options struct {
	DotFile    string
	ImgFile    string
	ImgFormat  string
	Stylesheet string
	YamlFile   string
}

//ParseArgs parses command line arguments into Options.
//It returns Options and a boolean flag indicating if
//help menu was requested.
func ParseArgs(args []string) (opts Options, h bool) {
	o := Options{
		DotFile:    "micromap.dot",
		ImgFile:    "micromap.png",
		ImgFormat:  "png",
		Stylesheet: "micromap.style",
		YamlFile:   "micromap.yml",
	}

	for _, arg := range args[1:] {
		optVal := strings.Split(arg, "=")
		switch optVal[0] {
		case "-d", "--dot":
			o.DotFile = optVal[1]
		case "-f", "--format":
			o.ImgFormat = optVal[1]
		case "-h", "--help":
			fmt.Print(micromapHelp)
			return o, true
		case "-i", "--img":
			o.ImgFile = optVal[1]
		case "-s", "--style":
			o.Stylesheet = optVal[1]
		case "-y", "--yaml":
			o.YamlFile = optVal[1]
		}
	}
	return o, false
}
