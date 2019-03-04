package micromap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

type Micromap struct {
	Config *Configuration
	Groups []Group
	Deps   []Dependency
	Rels   []Relation
}

type Configuration struct {
	App string
}

type Dependency struct {
	Name   string
	Typ    string
	Parent string
}

type Relation struct {
	Service string
	Path    string
	Dir     string
}
type Group struct {
	Name string
}

func MapMany(filenames []string) Micromap {
	var mmap Micromap
	for _, file := range filenames {
		m := Map(file)
		if m.Config.App != "" {
			mmap.Config = m.Config
		}
		mmap.Groups = append(mmap.Groups, m.Groups...)
		mmap.Deps = append(mmap.Deps, m.Deps...)
		mmap.Rels = append(mmap.Rels, m.Rels...)
	}
	return mmap
}

func Map(filename string) Micromap {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var cfg Configuration
	var grps []Group
	var deps []Dependency
	var rels []Relation

	regex, err := regexp.Compile(`@micromap.(config|group|dep|rel){[^}]*}`)

	finds := regex.FindAll(file, -1)
	for _, find := range finds {
		skip := 10
		var err error
		switch string(find[skip:13]) {
		case "con":
			content := find[skip+6:]
			err = json.Unmarshal(content, &cfg)
		case "gro":
			content := find[skip+5:]
			var grp Group
			err = json.Unmarshal(content, &grp)
			if err == nil {
				grps = append(grps, grp)
			}
		case "dep":
			content := find[skip+3:]
			var dep Dependency
			err = json.Unmarshal(content, &dep)
			if err == nil {
				deps = append(deps, dep)
			}
		case "rel":
			content := find[skip+3:]
			var rel Relation
			err = json.Unmarshal(content, &rel)
			if rel.Dir == "" {
				rel.Dir = "none"
			}
			if err == nil {
				rels = append(rels, rel)
			}
		}
		if err != nil {
			fmt.Println("missed marker:", string(find))
			log.Println(err)
		}
	}
	return Micromap{
		Config: &cfg,
		Groups: grps,
		Deps:   deps,
		Rels:   rels,
	}
}
