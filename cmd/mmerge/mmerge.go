package main

import (
	"fmt"
	"github.com/alvarg93/micromap/pkg/dot"
	"github.com/alvarg93/micromap/pkg/m4"
	"github.com/alvarg93/micromap/pkg/opts"
	"log"
	"os"
)

const intermediateFile = ".micromap.merged.dot"

func main() {
	opts, help := opts.ParseArgs(os.Args)
	if help {
		os.Exit(0)
	}

	mergeFile, err := os.Create(intermediateFile)
	if err != nil {
		log.Panic(err)
	}

	dot.Merge(mergeFile, opts.Sources)
	m4.Merge(intermediateFile, opts.DotFile)
	dot.ToPng(opts.DotFile, opts.ImgFile, opts.ImgFormat)

	mergeFile.Close()
	err = os.Remove(intermediateFile)
	if err != nil {
		fmt.Println(err)
	}
}
