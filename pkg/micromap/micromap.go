package micromap

import (
	"bufio"
	"gopkg.in/yaml.v2"
	"io"
)

// Micromap contains all relevant dependency data
type Micromap struct {
	Apps []App `yaml:"apps"`
}

// App represents a single service/program/task
type App struct {
	Name   string  `yaml:"name"`
	Groups []Group `yaml:"groups,omitempty"`
}

// Dependency represents the @micromap.dep annotation
type Dependency struct {
	Name   string     `yaml:"name"`
	Parent string     `yaml:"parent,omitempty"`
	Rels   []Relation `yaml:"relations,omitempty"`
	Type   string     `yaml:"type,omitempty"`
}

//Group represents the @micromap.group annotation
type Group struct {
	Deps []Dependency `yaml:"dependencies,omitempty"`
	Name string       `yaml:"name"`
}

// Relation represents the @micromap.rel annotation
type Relation struct {
	Dir     string `yaml:"dir,omitempty"`
	Name    string `yaml:"name"`
	Owner   string `yaml:"owner,omitempty"`
	Service string `yaml:"service,omitempty"`
}

// FromYaml takes an io.Reader and reads from it, expecting yaml syntax
func FromYaml(yml io.Reader) (Micromap, error) {
	scanner := bufio.NewScanner(yml)
	content := make([]byte, 0, 1024)
	for scanner.Scan() {
		content = append(content, scanner.Bytes()...)
		content = append(content, byte('\n'))
	}
	m := Micromap{}
	return m, yaml.Unmarshal(content, &m)
}

// Groups ...
func (m Micromap) Groups() []Group {
	groupsByName := map[string]Group{}
	for _, app := range m.Apps {
		for _, g := range app.Groups {
			mGrp, hasGrp := groupsByName[g.Name]
			if !hasGrp {
				groupsByName[g.Name] = g
				continue
			}
			mGrp.Deps = append(mGrp.Deps, g.Deps...)
			groupsByName[g.Name] = mGrp
		}
	}

	groups := []Group{}
	for _, g := range groupsByName {
		groups = append(groups, g)
	}
	return groups
}
