package opts_test

import (
	"github.com/lukaszjanyga/micromap/pkg/opts"
	"testing"
)

var defOpts = opts.Options{
	Regex:     ".+",
	Recursive: false,
	DotFile:   "micromap.dot",
	ImgFile:   "micromap.png",
	ImgFormat: "png",
}

func TestDefaultOpts(t *testing.T) {
	o, h := opts.ParseArgs([]string{})
	if h {
		t.Error("help triggered for default opts")
	}
	if o != defOpts {
		t.Error("options are different than default for empty args")
	}
}

func makeOpts(x, d, i, f string, r bool) opts.Options {
	o, _ := opts.ParseArgs([]string{})
	if x != "" {
		o.Regex = x
	}
	if d != "" {
		o.DotFile = d
	}
	if i != "" {
		o.ImgFile = i
	}
	if f != "" {
		o.ImgFormat = f
	}
	o.Recursive = r
	return o
}

var testArgs = []struct {
	args []string
	o    opts.Options
	h    bool
}{
	{[]string{"-h"}, defOpts, true},
	{[]string{"--help"}, defOpts, true},
	{[]string{"-r"}, makeOpts("", "", "", "", true), false},
	{[]string{"--recursive"}, makeOpts("", "", "", "", true), false},
	{[]string{"-x=regex"}, makeOpts("regex", "", "", "", false), false},
	{[]string{"--regex=regex"}, makeOpts("regex", "", "", "", false), false},
	{[]string{"-d=dot"}, makeOpts("", "dot", "", "", false), false},
	{[]string{"--dot=dot"}, makeOpts("", "dot", "", "", false), false},
	{[]string{"-i=img"}, makeOpts("", "", "img", "", false), false},
	{[]string{"--img=img"}, makeOpts("", "", "img", "", false), false},
	{[]string{"-f=format"}, makeOpts("", "", "", "format", false), false},
	{[]string{"--format=format"}, makeOpts("", "", "", "format", false), false},
	{[]string{"-r", "-x=regex", "-d=dot", "-i=img", "-f=format"}, makeOpts("regex", "dot", "img", "format", true), false},
}

func TestParseArgs(t *testing.T) {
	for _, testArg := range testArgs {
		o, h := opts.ParseArgs(testArg.args)
		if h != testArg.h {
			t.Error("help incorrect for", testArg.args)
		}
		if o != testArg.o {
			t.Error("options incorrect for", testArg.args)
		}
	}
}
