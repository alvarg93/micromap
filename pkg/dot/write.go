package dot

import (
	"fmt"
	"github.com/lukaszjanyga/micromap/pkg/color"
	"github.com/lukaszjanyga/micromap/pkg/micromap"
	"io"
)

const rankdir = "rankdir=LR"
const graphDefaults = `graph [splines="line", ranksep="2", nodesep="1"];` //throws errors on merges, when graphs become subgraphs
const nodeDefaults = `node [margin=0.25, fontcolor=white, style="filled,rounded", fontname = "sans-serif"];` + "\n"
const edgeDefaults = `edge [constraint=true, fontname = "sans-serif"];`
const nodeTemplate = `"%s"[shape=%s, fillcolor="%s", color="%s"]`
const edgeTemplate = `"%s" -- "%s"[dir=%s, headlabel="%s"];`

const s = 0.4
const v = 0.8

//Write writes the contents of a Micromap into a file.
func Write(f io.Writer, mmap micromap.Micromap) error {
	var content string

	content += beginGraph()
	subgraphs := 0

	for name, app := range mmap.Apps {
		fmt.Println(app)
		content += node(name, "", color.RandomHSV(s, v))

		grpDeps := make(map[string][]string)

		for _, grp := range app.Groups {
			for _, dep := range grp.Deps {
				content += node(dep.Name, dep.Typ, color.RandomHSV(s, v))
			}
		}
		for _, dep := range app.Deps {
			content += node(dep.Name, dep.Typ, color.RandomHSV(s, v))
			if dep.Parent != "" {
				grpDeps[dep.Parent] = append(grpDeps[dep.Parent], dep.Name)
			}
		}

		for _, grp := range app.Groups {
			for _, dep := range grp.Deps {
				for _, rel := range dep.Rels {
					content += edge(name, dep.Name, rel.Name, rel.Dir)
				}
			}
		}
		for _, dep := range app.Deps {
			for _, rel := range dep.Rels {
				content += edge(name, dep.Name, rel.Name, rel.Dir)
			}
		}
		for _, rel := range app.Rels {
			content += edge(name, rel.Service, rel.Name, rel.Dir)
		}

		for i, grp := range app.Groups {
			content += beginGroup(subgraphs+i, grp.Name)
			for _, dep := range grpDeps[grp.Name] {
				content += `"` + dep + `";`
			}
			for _, dep := range grp.Deps {
				content += `"` + dep.Name + `";`
			}
			content += endGroup()
		}
		subgraphs += len(app.Groups)
	}
	content += endGraph()

	_, err := f.Write([]byte(content))
	return err
}

func beginGraph() string {
	content := "graph {\n"
	content += rankdir + "\n"
	// content += graphDefaults + "\n" //throws errors on merges, when graphs become subgraphs
	content += nodeDefaults + "\n"
	content += edgeDefaults + "\n"
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

func node(name, typ string, color color.HSV) string {
	shape := "box"
	switch typ {
	case "db":
		shape = "cylinder"
	case "queue":
		shape = "box3d"
	case "lib":
		shape = "tab"
	}
	borderColor := color
	borderColor.V += 0.1
	return fmt.Sprintf(nodeTemplate, name, shape, color.ToRGB().ToColorHash(), borderColor.ToRGB().ToColorHash()) + "\n"
}
