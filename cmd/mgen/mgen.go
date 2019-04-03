package main

import (
	"fmt"
	"github.com/alvarg93/micromap/pkg/files"
	"github.com/alvarg93/micromap/pkg/micromap"
	"github.com/alvarg93/micromap/pkg/opts"
	myYaml "github.com/alvarg93/micromap/pkg/yaml"
	"log"
	"os"
)

func main() {
	opts, help := opts.ParseArgs(os.Args)
	if help {
		os.Exit(0)
	}

	fs, err := files.FindFiles(opts)
	if err != nil {
		log.Panic("could not open files", err)
	}

	mmap := micromap.MapManyYaml(fs)
	err = myYaml.SaveToFile(opts.YamlFile, mmap)
	if err != nil {
		fmt.Println(err)
	}
}
