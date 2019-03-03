package micromap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

type Entry struct {
	Config Configuration
	Dep    Dependency
}

type Configuration struct {
	App string
}

type Dependency struct {
	Channel  string
	Dir      string
	Typ      string
	Endpoint string
	Service  string
}

func IndexFile(path string) []Entry {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var entries []Entry
	regex, err := regexp.Compile(`@micromap[|][^|]*[|]`)

	finds := regex.FindAll(file, -1)
	for _, find := range finds {
		content := find[10 : len(find)-1]
		var e Entry
		err := json.Unmarshal(content, &e)
		if err == nil {
			entries = append(entries, e)
		} else {
			fmt.Println("missed marker:", string(content))
			log.Println(err)
		}
	}
	return entries
}
