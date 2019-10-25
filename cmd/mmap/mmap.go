package main

import (
	"fmt"
	"github.com/lukaszjanyga/micromap/pkg/dot"
	"github.com/lukaszjanyga/micromap/pkg/micromap"
	"github.com/lukaszjanyga/micromap/pkg/opts"
	"math/rand"
	"os"
	"time"
)

func main() {
	opts, help := opts.ParseArgs(os.Args)
	if help {
		os.Exit(0)
	}

	rand.Seed(time.Now().Unix())

	yml, err := os.Open(opts.YamlFile)
	if yml != nil {
		defer yml.Close()
	}
	if err != nil {
	}
	mmap := micromap.FromYaml(yml)

	f, err := os.Create(opts.DotFile)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = dot.Write(f, mmap)
	if err != nil {
		fmt.Println("failed to write dot file", err)
		os.Exit(1)
	}
	err = dot.ToPng(opts.DotFile, opts.ImgFile, opts.ImgFormat)
	if err != nil {
		fmt.Println("failed to write dot file", err)
		os.Exit(1)
	}
}
