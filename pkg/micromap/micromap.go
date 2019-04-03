package micromap

import (
	"bufio"
	"bytes"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"regexp"
)

//Micromap contains the result of parsing @micromap annotations
type Micromap struct {
	Config *Configuration `yaml:"config,omitempty"`
	Groups []Group        `yaml:"groups,omitempty"`
	Deps   []Dependency   `yaml:"deps,omitempty"`
	Rels   []Relation     `yaml:"rels,omitempty"`
}

//Configuration represents the @micromap.config annotation
type Configuration struct {
	App *string `yaml:"app,omitempty"`
}

//Dependency represents the @micromap.dep annotation
type Dependency struct {
	Name   string     `yaml:"name,omitempty"`
	Typ    string     `yaml:"typ,omitempty"`
	Parent string     `yaml:"parent,omitempty"`
	Rels   []Relation `yaml:"rels,omitempty"`
}

//Relation represents the @micromap.rel annotation
type Relation struct {
	Service string `yaml:"service,omitempty"`
	Path    string `yaml:"path,omitempty"`
	Dir     string `yaml:"dir,omitempty"`
}

//Group represents the @micromap.group annotation
type Group struct {
	Name string       `yaml:"name,omitempty"`
	Deps []Dependency `yaml:"deps,omitempty"`
}

func FromYaml(yml io.Reader) Micromap {
	scanner := bufio.NewScanner(yml)
	var content []byte
	for scanner.Scan() {
		content = append(content, scanner.Bytes()...)
		content = append(content, byte('\n'))
	}
	m := Micromap{}
	yaml.Unmarshal(content, &m)
	return m
}

//MapMany creates a Micromap containing values from
//all provided filenames.
func MapManyYaml(filenames []string) Micromap {
	var mmap Micromap
	for _, file := range filenames {
		m := MapYaml(file)
		if m.Config != nil {
			mmap.Config = m.Config
		}
		mmap.Groups = append(mmap.Groups, m.Groups...)
		mmap.Deps = append(mmap.Deps, m.Deps...)
		mmap.Rels = append(mmap.Rels, m.Rels...)
	}
	return mmap
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
		mmap := Micromap{
			Config: &Configuration{},
		}
		if len(result) >= 9 && string(result[:9]) == "@micromap" {
			results := bytes.SplitN(result, []byte{'\n'}, 2)
			err = yaml.Unmarshal(results[1], &mmap)
			if err != nil {
				log.Println("missed marker:", string(find))
				log.Println(err)
				continue
			}
		}
		if mmap.Config != nil && mmap.Config.App != nil {
			m.Config = mmap.Config
		}
		for _, g := range mmap.Groups {
			for _, d := range g.Deps {
				for _, r := range d.Rels {
					r.Service = d.Name
					m.Rels = append(m.Rels, r)
				}
				d.Rels = nil
				d.Parent = g.Name
				m.Deps = append(m.Deps, d)
			}
			g.Deps = nil
			m.Groups = append(m.Groups, g)
		}

		// m.Groups = append(m.Groups, mmap.Groups...)
		// m.Deps = append(m.Deps, mmap.Deps...)
		// m.Rels = append(m.Rels, mmap.Rels...)
	}
	return m
}

func trimInlineComments(content []byte) []byte {
	re := regexp.MustCompile(`\s+(//|#)+`)
	return re.ReplaceAll(content, []byte{'\n'})
}
