package dot

import (
	"fmt"
	"github.com/lukaszjanyga/micromap/pkg/color"
	"github.com/lukaszjanyga/micromap/pkg/micromap"
	"io"
)

const rankdir = "rankdir=LR\n"
const graphDefaults = "graph [dpi=600, splines=\"spline\", ranksep=\"2\", nodesep=\"1\"];\n"
const nodeDefaults = "node [margin=0.25, fontcolor=white, style=\"filled,rounded\", fontname = \"sans-serif\"];\n"
const edgeDefaults = "edge [constraint=true, fontname = \"sans-serif\"];\n"
const nodeTemplate = "\"%s\"[shape=%s, fillcolor=\"%s\", color=\"%s\"]\n"
const edgeTemplate = "\"%s\" -- \"%s\"[dir=%s, headlabel=\"%s\"];\n"

const s = 0.4
const v = 0.8

type Dot struct {
	*micromap.Micromap
	err error
}

//Write writes the contents of a Micromap into a file.
func (d Dot) Write(w io.Writer) error {
	d.beginGraph(w)
	subgraphs := 0

	for _, app := range d.Micromap.Apps {
		grpMap := map[string]struct{}{}
		depMap := map[string]struct{}{}
		// grpDeps := make(map[string][]string)

		for _, grp := range app.Groups {
			if _, exists := grpMap[grp.Name]; !exists {
				d.node(w, grp.Name, "group", color.HSV{H: 0, S: 0, V: 100})
				grpMap[grp.Name] = struct{}{}
			}
			for _, dep := range grp.Deps {
				if _, isAlsoGroup := grpMap[dep.Name]; isAlsoGroup {
					continue
				}
				if _, exists := depMap[dep.Name]; exists {
					continue
				}
				d.node(w, dep.Name, dep.Type, color.RandomHSV(s, v))
			}
		}

		// for _, dep := range app.Deps {
		// 	if _, exists := grpMap[dep.Parent]; !exists {
		// 		node(w, dep.Parent, "group", color.HSV{H: 0, S: 0, V: 100})
		// 		grpMap[dep.Parent] = struct{}{}
		// 	}

		// 	if _, exists := depMap[dep.Name]; !exists {
		// 		if _, isAlsoGroup := grpMap[dep.Name]; !isAlsoGroup {
		// 			node(w, dep.Name, dep.Type, color.RandomHSV(s, v))
		// 		}
		// 	}
		// 	if dep.Parent != "" {
		// 		grpDeps[dep.Parent] = append(grpDeps[dep.Parent], dep.Name)
		// 	}
		// }

		if _, exists := grpMap[app.Name]; !exists {
			d.node(w, app.Name, "", color.RandomHSV(s, v))
		}

		for _, grp := range app.Groups {
			for _, dep := range grp.Deps {
				for _, rel := range dep.Rels {
					d.edge(w, app.Name, dep.Name, rel.Name, rel.Dir)
				}
			}
		}
		// for _, dep := range app.Deps {
		// 	for _, rel := range dep.Rels {
		// 		edge(w, app.Name, dep.Name, rel.Name, rel.Dir)
		// 	}
		// }
		// for _, rel := range app.Rels {
		// 	edge(w, app.Name, rel.Service, rel.Name, rel.Dir)
		// }

	}

	for i, grp := range d.Micromap.Groups() {
		d.beginGroup(w, subgraphs+i, grp.Name)
		d.write(w, []byte(`"`+grp.Name+`";`))
		for _, dep := range grp.Deps {
			d.write(w, []byte(`"`+dep.Name+`";`))
		}
		d.endGroup(w)
	}

	d.endGraph(w)
	return d.err
}

func (d *Dot) beginGraph(w io.Writer) {
	d.write(w, []byte("graph {\n"))
	d.write(w, []byte(rankdir))
	d.write(w, []byte("layout=fdp;\n"))
	d.write(w, []byte("size=\"3,3\";\n"))
	d.write(w, []byte("overlap=false;\n"))
	d.write(w, []byte("splines=false;\n"))
	// d.write(w, []byte("pack=true;\n"))
	d.write(w, []byte("start=\"random\";\n"))
	d.write(w, []byte("sep=0.8;\n"))
	d.write(w, []byte("inputscale=0.4;\n"))
	// K=0.50
	// maxiter=2000
	// start=1251

	d.write(w, []byte("K=0.50;\n"))
	d.write(w, []byte("maxiter=2000;\n"))
	d.write(w, []byte("start=1251;\n"))
	d.write(w, []byte(graphDefaults))
	d.write(w, []byte(nodeDefaults))
	d.write(w, []byte(edgeDefaults))
}

func (d *Dot) endGraph(w io.Writer) {
	d.write(w, []byte("}\n"))
}

func (d *Dot) beginGroup(w io.Writer, number int, name string) {
	d.write(w, []byte(fmt.Sprintf("subgraph cluster_%d{\n", number)))
	d.write(w, []byte(fmt.Sprintf("label=\"%s\";\n", name)))
}

func (d *Dot) endGroup(w io.Writer) {
	d.write(w, []byte("\n}\n"))
}

func (d *Dot) edge(w io.Writer, a, b, label, dir string) {
	if dir == "" {
		dir = "none"
	}
	d.write(w, []byte(fmt.Sprintf(edgeTemplate, a, b, dir, label)))
}

func (d *Dot) node(w io.Writer, name, typ string, color color.HSV) {
	shape := "box"
	switch typ {
	case "db":
		shape = "cylinder"
	case "queue":
		shape = "box3d"
	case "lib":
		shape = "tab"
	case "group":
		shape = "point"
	}
	borderColor := color
	borderColor.V += 0.1
	d.write(w, []byte(fmt.Sprintf(nodeTemplate, name, shape, color.ToRGB().ToColorHash(), borderColor.ToRGB().ToColorHash())))
}

func (d *Dot) write(w io.Writer, p []byte) {
	_, err := w.Write(p)
	if err != nil {
		fmt.Println(err.Error())
	}
	if d.err != nil {
		d.err = err
	}
}
