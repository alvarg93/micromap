package main

import (
	"github.com/alvarg93/micromap/pkg/dot"
	"github.com/alvarg93/micromap/pkg/files"
	"github.com/alvarg93/micromap/pkg/micromap"
	"github.com/alvarg93/micromap/pkg/opts"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	opts, help := opts.ParseArgs(os.Args[1:])
	if help {
		os.Exit(0)
	}

	rand.Seed(time.Now().Unix())

	fs, err := files.FindFiles(opts)
	if err != nil {
		log.Panic("could not open files", err)
	}

	mmap := micromap.MapMany(fs)

	f, err := os.Create(opts.DotFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = dot.WriteToFile(f, mmap)
	if err != nil {
		log.Panic("failed to write dot file", err)
	}
	err = dot.ToPng(opts.DotFile, opts.ImgFile, opts.ImgFormat)
	if err != nil {
		log.Panic("failed to write img file", err)
	}
}
