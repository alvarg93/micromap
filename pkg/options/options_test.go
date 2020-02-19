package options_test

import (
	"github.com/lukaszjanyga/micromap/pkg/opts"
	"testing"
)

var defOpts = opts.Options{
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

func makeOpts(d, i, f string) opts.Options {
	o, _ := opts.ParseArgs([]string{})
	if d != "" {
		o.DotFile = d
	}
	if i != "" {
		o.ImgFile = i
	}
	if f != "" {
		o.ImgFormat = f
	}
	return o
}

var testArgs = []struct {
	args []string
	o    opts.Options
	h    bool
}{
	{[]string{"-h"}, defOpts, true},
	{[]string{"--help"}, defOpts, true},
	{[]string{"-r"}, makeOpts("", "", ""), false},
	{[]string{"--recursive"}, makeOpts("", "", ""), false},
	{[]string{"-x=regex"}, makeOpts("", "", ""), false},
	{[]string{"--regex=regex"}, makeOpts("", "", ""), false},
	{[]string{"-d=dot"}, makeOpts("dot", "", ""), false},
	{[]string{"--dot=dot"}, makeOpts("dot", "", ""), false},
	{[]string{"-i=img"}, makeOpts("", "img", ""), false},
	{[]string{"--img=img"}, makeOpts("", "img", ""), false},
	{[]string{"-f=format"}, makeOpts("", "", "format"), false},
	{[]string{"--format=format"}, makeOpts("", "", "format"), false},
	{[]string{"-r", "-x=regex", "-d=dot", "-i=img", "-f=format"}, makeOpts("dot", "img", "format"), false},
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
