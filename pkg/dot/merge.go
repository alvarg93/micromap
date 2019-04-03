package dot

import (
	"fmt"
	"github.com/alvarg93/micromap/pkg/opts"
	"io"
)

func Merge(writer io.Writer, sources []opts.GraphSource) (int, error) {
	var content string
	content += "graph merged {\n"
	content += "rankdir=LR\n"
	content += "define(graph,subgraph)\n"

	for i, src := range sources {
		fmt.Println(i, src)
		if !src.IsURL {
			content += fmt.Sprintf("include(%s)\n", src.Path)
		}
	}

	content += "}\n"

	return writer.Write([]byte(content))
}
