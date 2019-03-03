package main

import (
	"encoding/hex"
	"fmt"
	"github.com/alvarg93/micromap/pkg/files"
	"github.com/alvarg93/micromap/pkg/micromap"
	"log"
	"math/rand"
	"os"
	"os/exec"
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

	f, err := os.Create("micromap.dot")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString("graph {\n")
	_, err = f.WriteString("\"" + config.App + "\"[fontcolor=white,style=filled,color=\"" + getColorHash() + "\"]\n")
	for service, rels := range deps {
		_, err = f.WriteString("\"" + service + "\"[fontcolor=white,style=filled,color=\"" + getColorHash() + "\"]\n")
		for _, rel := range rels {
			_, err = f.WriteString(edge(config.App, service, rel.Endpoint+rel.Channel, "1") + "\n")
		}
	}
	_, err = f.WriteString("}\n")

	_, err = exec.Command("sh", "-c", "dot -Tpng micromap.dot -o micromap.png").Output()
}

func edge(a, b, l, w string) string {
	return fmt.Sprintf("\"%s\" -- \"%s\"[label=\"%s\",weight=\"%s\"];", a, b, l, w)
}

func getColorHash() string {
	r, g, b := uint8(rand.Uint32()%200), uint8(rand.Uint32()%200), uint8(rand.Uint32()%200)
	return "#" + hex.EncodeToString([]byte{byte(r), byte(g), byte(b)})
}
