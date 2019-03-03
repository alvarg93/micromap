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

func main() {
	fs, err := files.FindFiles(os.Args[1])
	if err != nil {
		log.Fatal("Could not open files")
	}

	var mmap micromap.Micromap
	for _, file := range fs {
		m := micromap.Map(file)
		if m.Config.App != "" {
			mmap.Config = m.Config
		}
		mmap.Groups = append(mmap.Groups, m.Groups...)
		mmap.Deps = append(mmap.Deps, m.Deps...)
		mmap.Rels = append(mmap.Rels, m.Rels...)
	}

	f, err := os.Create("micromap.dot")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var content string

	content += "graph {\n"
	content += "rankdir=LR\n"
	content += node(mmap.Config.App, getColorHash(), "circle") + "\n"

	grpDeps := make(map[string][]string)

	for _, dep := range mmap.Deps {
		content += node(dep.Name, getColorHash(), "box3d") + "\n"
		if dep.Parent != "" {
			grpDeps[dep.Parent] = append(grpDeps[dep.Parent], dep.Name)
		}
	}

	for _, rel := range mmap.Rels {
		content += edge(mmap.Config.App, rel.Service, rel.Dir+":"+rel.Path, "1") + "\n"
	}

	for i, grp := range mmap.Groups {
		content += fmt.Sprintf("subgraph cluster_%d{\n", i)
		content += "label=\"" + grp.Name + "\";\n"
		for _, dep := range grpDeps[grp.Name] {
			content += dep + ";"
		}
		content += "\n}\n"
	}
	content += "}\n"
	if err != nil {
		panic(err)
	}

	_, err = f.WriteString(content)
	if err != nil {
		panic(err)
	}

	_, err = exec.Command("sh", "-c", "dot -Tpng micromap.dot -o micromap.png").Output()
	if err != nil {
		panic(err)
	}

}

func edge(a, b, label, weight string) string {
	return fmt.Sprintf("\"%s\" -- \"%s\"[label=\"%s\",weight=\"%s\"];", a, b, label, weight)
}

func node(name, color, shape string) string {
	return "\"" + name + "\"[shape=" + shape + ",fontcolor=white,style=filled,fillcolor=\"" + color + "\"]"
}

func getColorHash() string {
	r, g, b := uint8(rand.Uint32()%200), uint8(rand.Uint32()%200), uint8(rand.Uint32()%200)
	return "#" + hex.EncodeToString([]byte{byte(r), byte(g), byte(b)})
}
