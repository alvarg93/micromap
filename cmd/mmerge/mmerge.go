package main

import (
	"fmt"
	"github.com/lukaszjanyga/micromap/pkg/dot"
	"github.com/lukaszjanyga/micromap/pkg/m4"
	"github.com/lukaszjanyga/micromap/pkg/opts"
	"os"
)

const intermediateFile = ".micromap.merged.dot"

func main() {
	var err error
	var tempFiles []string
	defer func() {
		fmt.Println("Removing temp files...")
		var errors []error
		// err = os.Remove(intermediateFile)
		// if err != nil {
		// 	errors = append(errors, err)
		// }
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
	}()

	opts, help := opts.ParseArgs(os.Args)
	if help {
		return
	}

	fmt.Println("Creating temp files...")
	mergeFile, err := os.Create(intermediateFile)
	if mergeFile != nil {
		defer mergeFile.Close()
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Merging dot files...")
	tempFiles, err = dot.Merge(mergeFile, opts.Sources)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Merging...")
	err = m4.Merge(intermediateFile, opts.DotFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Converting to PNG...")
	err = dot.ToPng(opts.DotFile, opts.ImgFile, opts.ImgFormat)
	if err != nil {
		fmt.Println(err)
		return
	}
}
