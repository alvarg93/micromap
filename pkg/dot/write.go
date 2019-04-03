package dot

import (
	"encoding/hex"
	"fmt"
	"github.com/lukaszjanyga/micromap/pkg/micromap"
	"io"
	"math/rand"
)

const nodeTemplate = `"%s"[shape=%s,fontcolor=white,style=filled,fillcolor="%s"]`
const edgeTemplate = `"%s" -- "%s"[dir=%s,label="%s"];`

//Write writes the contents of a Micromap into a file.
func Write(f io.Writer, mmap micromap.Micromap) error {
	var content string

	content += beginGraph()
	if mmap.Config != nil {
		content += node(*mmap.Config.App, getColorHash(), "")
	}

	grpDeps := make(map[string][]string)

	for _, grp := range mmap.Groups {
		for _, dep := range grp.Deps {
			content += node(dep.Name, getColorHash(), dep.Typ)
		}
	}
	for _, dep := range mmap.Deps {
		content += node(dep.Name, getColorHash(), dep.Typ)
		if dep.Parent != "" {
			grpDeps[dep.Parent] = append(grpDeps[dep.Parent], dep.Name)
		}
	}

	for _, grp := range mmap.Groups {
		for _, dep := range grp.Deps {
			for _, rel := range dep.Rels {
				content += edge(*mmap.Config.App, dep.Name, rel.Path, rel.Dir)
			}
		}
	}
	for _, dep := range mmap.Deps {
		for _, rel := range dep.Rels {
			content += edge(*mmap.Config.App, dep.Name, rel.Path, rel.Dir)
		}
	}
	for _, rel := range mmap.Rels {
		content += edge(*mmap.Config.App, rel.Service, rel.Path, rel.Dir)
	}

	for i, grp := range mmap.Groups {
		content += beginGroup(i, grp.Name)
		for _, dep := range grpDeps[grp.Name] {
			content += dep + ";"
		}
		for _, dep := range grp.Deps {
			content += dep.Name + ";"
		}
		content += endGroup()
	}
	content += endGraph()

	_, err := f.Write([]byte(content))
	return err
}

func beginGraph() string {
	content := "graph {\n"
	content += "rankdir=LR\n"
	return content
}

func endGraph() string {
	return "}\n"
}

func beginGroup(number int, name string) string {
	content := fmt.Sprintf("subgraph cluster_%d{\n", number)
	content += fmt.Sprintf(`label="%s";`, name) + "\n"
	return content
}

func endGroup() string {
	return "\n}\n"
}

func edge(a, b, label, dir string) string {
	if dir == "" {
		dir = "none"
	}
	return fmt.Sprintf(edgeTemplate, a, b, dir, label) + "\n"
}

func node(name, color, typ string) string {
	shape := "circle"
	switch typ {
	case "db":
		shape = "cylinder"
	case "queue":
		shape = "box3d"
	}
	return fmt.Sprintf(nodeTemplate, name, shape, color) + "\n"
}

func getColorHash() string {
	r, g, b := uint8(rand.Uint32()%127+100), uint8(rand.Uint32()%127+100), uint8(rand.Uint32()%127+100)
	return "#" + hex.EncodeToString([]byte{byte(r), byte(g), byte(b)})
}
