package micromap

import (
	"bufio"
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"regexp"
)

//Micromap contains the result of parsing @micromap annotations
type Micromap struct {
	Apps map[string]App `yaml:"apps"`
}

//Config represents the @micromap.config annotation
type App struct {
	// Name   string       `yaml:"name,omitempty"`
	Groups []Group      `yaml:"groups,omitempty"`
	Deps   []Dependency `yaml:"dependencies,omitempty"`
	Rels   []Relation   `yaml:"relations,omitempty"`
}

//Dependency represents the @micromap.dep annotation
type Dependency struct {
	Name   string     `yaml:"name,omitempty"`
	Typ    string     `yaml:"typ,omitempty"`
	Parent string     `yaml:"parent,omitempty"`
	Rels   []Relation `yaml:"relations,omitempty"`
}

//Relation represents the @micromap.rel annotation
type Relation struct {
	Name    string `yaml:"name,omitempty"`
	Service string `yaml:"service,omitempty"`
	Owner   string `yaml:"owner,omitempty"`
	Dir     string `yaml:"dir,omitempty"`
}

//Group represents the @micromap.group annotation
type Group struct {
	Name string       `yaml:"name,omitempty"`
	Deps []Dependency `yaml:"dependencies,omitempty"`
}

func FromYaml(yml io.Reader) Micromap {
	scanner := bufio.NewScanner(yml)
	var content []byte
	for scanner.Scan() {
		content = append(content, scanner.Bytes()...)
		content = append(content, byte('\n'))
	}
	m := Micromap{}
	err := yaml.Unmarshal(content, &m)
	if err != nil {
		fmt.Println(err)
	}
	return m
}

//MapMany creates a Micromap containing values from
//all provided filenames.
func MapManyYaml(filenames []string) Micromap {
	var m Micromap
	for _, file := range filenames {
		mmap := MapYaml(file)
		for key, val := range mmap.Apps {
			m.Apps[key] = val
		}
	}
	return m
}

//Map creates a Micromap containing values from
//the specified filename.
func MapYaml(filename string) Micromap {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	m := Micromap{}

	blockCmts := `\/\*([\s\S]*?)\*\/`
	inlnCmts := `(\s*(\/\/|#)\s*[^\r\n]+)+`
	regex := regexp.MustCompile(`(` + blockCmts + `|` + inlnCmts + `)`)

	finds := regex.FindAll(file, -1)
	for _, find := range finds {
		var err error
		result := trimInlineComments(find)[1:]
		mmap := Micromap{}
		if len(result) >= 9 && string(result[:9]) == "@micromap" {
			results := bytes.SplitN(result, []byte{'\n'}, 2)
			err = yaml.Unmarshal(results[1], &mmap)
			if err != nil {
				log.Println("missed marker:", string(find))
				log.Println(err)
				continue
			}
		}

		// for _, a := range mmap.Apps {
		// for _, g := range a.Groups {
		// 	for _, d := range g.Deps {
		// 		for _, r := range d.Rels {
		// 			r.Service = d.Name
		// 			m.Rels = append(m.Rels, r)
		// 		}
		// 		d.Rels = nil
		// 		d.Parent = g.Name
		// 		m.Deps = append(m.Deps, d)
		// 	}
		// 	g.Deps = nil
		// 	m.Groups = append(m.Groups, g)
		// }
		// for _, d := range mmap.Deps {
		// 	for _, r := range d.Rels {
		// 		r.Service = d.Name
		// 		m.Rels = append(m.Rels, r)
		// 	}
		// 	d.Rels = nil
		// 	m.Deps = append(m.Deps, d)
		// }
		// m.Rels = append(m.Rels, mmap.Rels...)
		// }
	}
	return m
}

func trimInlineComments(content []byte) []byte {
	re := regexp.MustCompile(`\s+(//|#)+`)
	return re.ReplaceAll(content, []byte{'\n'})
}
