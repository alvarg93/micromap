package main

import (
	"fmt"
	"github.com/lukaszjanyga/micromap/pkg/dot"
	"github.com/lukaszjanyga/micromap/pkg/micromap"
	"github.com/lukaszjanyga/micromap/pkg/options"
	"github.com/lukaszjanyga/micromap/pkg/png"
	"math/rand"
	"os"
	"time"
)

func main() {
	opts, help := options.ParseArgs(os.Args)
	if help {
		os.Exit(0)
	}

	rand.Seed(time.Now().Unix())

	yml, err := os.Open(opts.YamlFile)
	if yml != nil {
		defer yml.Close()
	}
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	mmap, err := micromap.FromYaml(yml)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	f, err := os.Create(opts.DotFile)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	err = dot.Dot{Micromap: &mmap}.Write(f)
	if err != nil {
		fmt.Println("failed to write dot file", err.Error())
		return
	}
	err = png.ToPng(opts.DotFile, opts.ImgFile, opts.ImgFormat)
	if err != nil {
		fmt.Println("failed to write dot file", err.Error())
		return
	}
}
