package dot

import (
	"encoding/hex"
	"fmt"
	"github.com/alvarg93/micromap/pkg/micromap"
	"math/rand"
	"os"
	"os/exec"
)

//WriteToFile writes the contents of a Micromap into a file.
func WriteToFile(f *os.File, mmap micromap.Micromap) error {
	var content string

	content += "graph {\n"
	content += "rankdir=LR\n"
	content += node(mmap.Config.App, getColorHash(), "app") + "\n"

	grpDeps := make(map[string][]string)

	for _, dep := range mmap.Deps {
		content += node(dep.Name, getColorHash(), dep.Typ) + "\n"
		if dep.Parent != "" {
			grpDeps[dep.Parent] = append(grpDeps[dep.Parent], dep.Name)
		}
	}

	for _, rel := range mmap.Rels {
		content += edge(mmap.Config.App, rel.Service, rel.Dir+":"+rel.Path, "1", rel.Dir) + "\n"
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

	_, err := f.WriteString(content)
	return err
}

//ToPng runs the dot command to convert a dot file into an image file
func ToPng(from, to, format string) error {
	_, err := exec.Command("sh", "-c", "dot -T"+format+" "+from+" -o "+to).Output()
	return err
}

func edge(a, b, label, weight, dir string) string {
	return fmt.Sprintf("\"%s\" -- \"%s\"[dir="+dir+",label=\"%s\"];", a, b, label)
}

func node(name, color, typ string) string {
	shape := "circle"
	switch typ {
	case "db":
		shape = "cylinder"
	case "queue":
		shape = "box3d"
	}
	return "\"" + name + "\"[shape=" + shape + ",fontcolor=white,style=filled,fillcolor=\"" + color + "\"]"
}

func getColorHash() string {
	r, g, b := uint8(rand.Uint32()%200), uint8(rand.Uint32()%200), uint8(rand.Uint32()%200)
	return "#" + hex.EncodeToString([]byte{byte(r), byte(g), byte(b)})
}
