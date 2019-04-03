package main

import (
	"fmt"
	"github.com/lukaszjanyga/micromap/pkg/files"
	"github.com/lukaszjanyga/micromap/pkg/micromap"
	"github.com/lukaszjanyga/micromap/pkg/opts"
	myYaml "github.com/lukaszjanyga/micromap/pkg/yaml"
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
