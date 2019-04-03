package yaml

import (
	"github.com/lukaszjanyga/micromap/pkg/micromap"
	"gopkg.in/yaml.v2"
	"os"
)

//SaveToFile saves a Micromap to a yaml file
func SaveToFile(filename string, mmap micromap.Micromap) error {
	result, err := yaml.Marshal(mmap)
	if err != nil {
		return err
	}

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(result)

	return err
}
