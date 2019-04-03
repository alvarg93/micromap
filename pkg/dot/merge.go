package dot

import (
	"fmt"
	"github.com/lukaszjanyga/micromap/pkg/http"
	"github.com/lukaszjanyga/micromap/pkg/opts"
	"io"
	"os"
)

func Merge(writer io.Writer, sources []opts.GraphSource) ([]string, error) {
	var content string
	content += "graph merged {\n"
	content += "rankdir=LR\n"
	content += "define(graph,subgraph)\n"

	tempFiles := []string{}
	for i, src := range sources {
		var path string
		fmt.Println("including:", src.Path)
		if !src.IsURL {
			path = src.Path
		} else {
			path := fmt.Sprintf("_micromap_temp_%d.dot", i)
			f, err := os.Create(path)
			if err != nil {
				return tempFiles, err
			}
			err = http.Download(f, src.Path)
			if err != nil {
				return tempFiles, err
			}
			tempFiles = append(tempFiles, path)
		}
		content += fmt.Sprintf("include(%s)\n", path)
	}

	content += "}\n"

	_, err := writer.Write([]byte(content))

	return tempFiles, err
}
