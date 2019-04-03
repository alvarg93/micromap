package main

import (
	"fmt"
	"github.com/lukaszjanyga/micromap/pkg/dot"
	"github.com/lukaszjanyga/micromap/pkg/m4"
	"github.com/lukaszjanyga/micromap/pkg/opts"
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

	tempFiles, err := dot.Merge(mergeFile, opts.Sources)
	if err != nil {
		errorExit(err)
	}

	err = m4.Merge(intermediateFile, opts.DotFile)
	if err != nil {
		errorExit(err)
	}

	err = dot.ToPng(opts.DotFile, opts.ImgFile, opts.ImgFormat)
	if err != nil {
		errorExit(err)
	}

	mergeFile.Close()

	var errors []error
	err = os.Remove(intermediateFile)
	if err != nil {
		errors = append(errors, err)
	}
	for _, tempFile := range tempFiles {
		err = os.Remove(tempFile)
		if err != nil {
			errors = append(errors, err)
		}
	}
	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Println(err)
		}
		os.Exit(1)
	}
}

func errorExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}
