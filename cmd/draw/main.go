package main

import (
	"fmt"
	"github.com/alvarg93/micromap/pkg/files"
	"github.com/alvarg93/micromap/pkg/micromap"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

/*
@micromap|{
	"config": {
		"app":"NOT YOUR MOM"
	},
	"dep": {
		"endpoint":"GET /api/v1/deliveries",
		"service":"your mom v1",
		"typ":"rest"
	}
}|
@micromap|{
	"dep": {
		"endpoint":"GET /api/v2/deliveries",
		"service":"your mom v2",
		"typ":"rest"
	}
}|
@micromap|{
	"dep": {
		"endpoint":"GET /api/v3/deliveries",
		"service":"your mom v3",
		"typ":"rest"
	}
}|
@micromap|{
	"dep": {
		"endpoint":"GET /api/v4/deliveries",
		"service":"your mom v4",
		"typ":"rest"
	}
}|
@micromap|{
	"dep": {
		"endpoint":"GET /api/v4/orders",
		"service":"your mom v4",
		"typ":"rest"
	}
}|
@micromap|{
	"dep": {
		"endpoint":"GET /api/v4/orders",
		"service":"your mom v5",
		"typ":"rest"
	}
}|
@micromap|{
	"dep": {
		"endpoint":"GET /api/v4/orders",
		"service":"your mom v6",
		"typ":"rest"
	}
}|
@micromap|{
	"dep": {
		"endpoint":"GET /api/v4/orders",
		"service":"your mom v7",
		"typ":"rest"
	}
}|
micromap|{
	"dep": {
		"endpoint":"GET /api/v4/orders",
		"service":"your mom v8",
		"typ":"rest"
	}
}|
micromap|{
	"dep": {
		"endpoint":"GET /api/v4/orders",
		"service":"your mom v9",
		"typ":"rest"
	}
}|
*/

/*
@micromap|{
	"dep": {
		"channel":"sync request",
		"dir":"out",
		"typ":"queue",
		"service":"sns"
	}
}|
*/

func main() {
	fs, err := files.FindFiles(os.Args[1])
	if err != nil {
		log.Fatal("Could not open files")
	}
	var cfgs []micromap.Configuration
	var deps map[string][]micromap.Dependency
	deps = make(map[string][]micromap.Dependency)
	for _, file := range fs {
		entries := micromap.IndexFile(file)
		for _, entry := range entries {
			cfgs = append(cfgs, entry.Config)
			deps[entry.Dep.Service] = append(deps[entry.Dep.Service], entry.Dep)
		}
	}

	var config micromap.Configuration
	for _, cfg := range cfgs {
		if cfg.App != "" {
			config = cfg
			break
		}
	}

	// TODO: WIP built in drawing library
	blockSize := 128
	max := cols(len(deps))
	size := (2*max - 1) * blockSize
	bcg := image.NewRGBA(image.Rect(0, 0, size, size))
	white := color.RGBA{255, 255, 255, 0}
	draw.Draw(bcg, bcg.Bounds(), &image.Uniform{white}, image.ZP, draw.Src)
	rand.Seed(time.Now().Unix())

	var srvNo int
	var row, col int
	for service, rels := range deps {
		if col == max/2 && row == max/2 {
			col++
		}
		for i, rel := range rels {
			fmt.Println(service, i, rel.Endpoint+rel.Channel)
			paintLine(bcg, blockSize*2*col+blockSize/2, blockSize*2*row+blockSize/2, size/2, size/2)
		}

		paintService(bcg, blockSize, col, row, service)

		if col+1 >= max {
			row++
			col = 0
		} else {
			col++
		}
		srvNo++
	}
	paintService(bcg, blockSize, max/2, max/2, config.App)

	f, err := os.Create("draw.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, bcg)
}

func paintService(bcg *image.RGBA, blockSize, col, row int, service string) {
	c := getColor()
	pt := image.Point{X: blockSize * 2 * col, Y: blockSize * 2 * row}
	m := image.NewRGBA(image.Rect(pt.X, pt.Y, pt.X+blockSize, pt.Y+blockSize))
	draw.Draw(bcg, m.Bounds(), &image.Uniform{c}, pt, draw.Src)
	addLabel(bcg, pt.X+10, pt.Y+blockSize/2, service)
}

func paintLine(bcg *image.RGBA, x0, y0, x1, y1 int) {
	println(x0, y0, x1, y1)
	hG := float64(x1 - x0)
	vG := float64(y1 - y0)
	blockSize := 8.0
	c := color.RGBA{200, 200, 255, 255}
	pt := image.Point{X: x0, Y: y0}
	for math.Abs(float64(pt.X-x1))+math.Abs(float64(pt.Y-y1)) > blockSize {
		m := image.NewRGBA(image.Rect(pt.X, pt.Y, pt.X+int(blockSize), pt.Y+int(blockSize)))
		draw.Draw(bcg, m.Bounds(), &image.Uniform{c}, pt, draw.Src)
		pt.X += int(hG / (10 * blockSize))
		pt.Y += int(vG / (10 * blockSize))
	}
}

func cols(n int) int {
	if n < 9 {
		return 3
	}
	return int(math.Ceil(math.Sqrt(float64(n))))
}

func getColor() color.Color {
	r, g, b := uint8(rand.Uint32()%128), uint8(rand.Uint32()%128), uint8(rand.Uint32()%128)
	return color.RGBA{r, g, b, 255}
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{255, 255, 255, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
