package files

import (
	"io/ioutil"
	"log"
	// "os"
	"path/filepath"
	"regexp"
)

func FindFiles(regex string) ([]string, error) {
	files, err := findInDir(".", regex)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func findInDir(path string, regex string) ([]string, error) {
	result := make([]string, 0)
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range dir {
		if ok, err := regexp.MatchString(regex, f.Name()); ok && err == nil && !f.IsDir() {
			p := filepath.Join(path, f.Name())
			abs, err := filepath.Abs(p)
			if err != nil {
				log.Println("failed to create absolute path for", p, "/", f.Name())
			}
			result = append(result, abs)
		} else if f.IsDir() {
			p := filepath.Join(path, f.Name())
			subdir, err := findInDir(p, regex)
			if err != nil {
				log.Println("failed to read from")
			}
			result = append(result, subdir...)
		}
	}
	return result, nil
}
