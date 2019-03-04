package files

import (
	"github.com/alvarg93/micromap/pkg/opts"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

//FindFiles scans the current directory for files matching
//opts.Regex and continues recursively if opts.Recursive is true.
func FindFiles(opts opts.Options) ([]string, error) {
	files, err := findInDir(".", opts.Regex, opts.Recursive)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func findInDir(path string, regex string, recursive bool) ([]string, error) {
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
		} else if f.IsDir() && recursive {
			p := filepath.Join(path, f.Name())
			subdir, err := findInDir(p, regex, recursive)
			if err != nil {
				log.Println("failed to read from")
			}
			result = append(result, subdir...)
		}
	}
	return result, nil
}
