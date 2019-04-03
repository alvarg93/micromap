package opts

import (
	"fmt"
	"path/filepath"
	"strings"
)

var micromapHelp = `
Micromap generates a dependency graph from a specification file in yaml.
Dependencies: dot (graphviz package)
Usage: mmap [-x=val][-r][-d=val][-i=val][-f=val]

	-h, --help        - this message
	-y, --yaml        - input file (default: micromap.yml)
	-d, --dot         - output file for a dot graph (default: micromap.dot)
	-i, --img         - output file for a visual representation of the graph (default: micromap.png)
	-f, --format      - format of the image file (default: png)
`

var mergeHelp = `
Micromerge pulls dependency specification from multiple dot files and merges them into one.
It temporarily downloads files from online sources.
Dependencies: m4 (graphviz package)
Usage: mmerge [-d=val][-i=val][-f=val] [list of paths or urls]

	-d, --dot         - output file for a dot graph (default: micromap.dot)
	-i, --img         - output file for a visual representation of the graph (default: micromap.png)
	-f, --format      - format of the image file (default: png)
`

var generateHelp = `
Microgen generates a specification file in yaml from code annotations.
Usage: mgen [-x=val][-r][-d=val][-i=val][-f=val]

	-h, --help        - this message
	-x, --regex       - regex pattern for files to scan (default: .)
	-r, --recursive   - scan recursive (default: false)
	-y, --yaml        - output file (default: micromap.yml)
`

//Options represent command line arguments for the program
type Options struct {
	Command   string
	DotFile   string
	ImgFile   string
	ImgFormat string
	Recursive bool
	Regex     string
	Sources   []GraphSource
	YamlFile  string
}

type GraphSource struct {
	Path     string
	FileName string
	IsURL    bool
}

var protocols = map[string]struct{}{"http": struct{}{}, "https": struct{}{}}

//ParseArgs parses command line arguments into Options.
//It returns Options and a boolean flag indicating if
//help menu was requested.
func ParseArgs(args []string) (opts Options, h bool) {
	o := Options{
		Command:   args[0],
		DotFile:   "micromap.dot",
		ImgFile:   "micromap.png",
		ImgFormat: "png",
		Recursive: false,
		Regex:     ".",
		YamlFile:  "micromap.yml",
	}

	args = args[1:]

	for _, arg := range args {
		optVal := strings.Split(arg, "=")
		switch optVal[0] {
		case "-h", "--help":
			printHelp(o.Command)
			return o, true
		case "-d", "--dot":
			o.DotFile = optVal[1]
		case "-i", "--img":
			o.ImgFile = optVal[1]
		case "-f", "--format":
			o.ImgFormat = optVal[1]
		case "-r", "--recursive":
			o.Recursive = true
		case "-x", "--regex":
			o.Regex = optVal[1]
		case "-y", "--yaml":
			o.YamlFile = optVal[1]
		default:
			o.Sources = append(o.Sources, parseSource(arg))
		}
	}
	return o, false
}

//printHelp prints help specific for command
func printHelp(command string) {
	switch command {
	case "mmap", "./mmap":
		fmt.Print(micromapHelp)
	case "mmerge", "./mmerge":
		fmt.Print(mergeHelp)
	case "mgen", "./mgen":
		fmt.Print(generateHelp)
	}
}

//parseSource parses an argument into a GraphSource
func parseSource(source string) GraphSource {
	gs := GraphSource{}
	abs, err := filepath.Abs(source)
	if err != nil {
		return gs
	}

	protocolSplit := strings.Split(source, "://")
	if len(protocolSplit) == 2 {
		_, isProtocol := protocols[protocolSplit[0]]
		gs.IsURL = isProtocol
	}
	gs.FileName = filepath.Base(abs)
	gs.Path = abs
	return gs
}
