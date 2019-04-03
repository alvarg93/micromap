package main

import (
	"github.com/alvarg93/micromap/pkg/dot"
	"github.com/alvarg93/micromap/pkg/micromap"
	"github.com/alvarg93/micromap/pkg/opts"
	"log"
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

	yml, _ := os.Open(opts.YamlFile)
	defer yml.Close()
	mmap := micromap.FromYaml(yml)

	f, err := os.Create(opts.DotFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = dot.Write(f, mmap)
	if err != nil {
		log.Panic("failed to write dot file", err)
	}
	err = dot.ToPng(opts.DotFile, opts.ImgFile, opts.ImgFormat)
	if err != nil {
		log.Panic("failed to write img file", err)
	}
}
