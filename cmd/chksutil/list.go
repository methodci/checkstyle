package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/subcommands"
	"github.com/methodci/checkstyle"
)

type ListCmd struct {
	maxLineShift int
}

func (*ListCmd) Name() string     { return "list" }
func (*ListCmd) Synopsis() string { return "list contents of one or more checkstyle files" }
func (*ListCmd) Usage() string {
	return `list <left-file> <right-file>:
	list contents of one or more checkstyle files.
`
}

func (p *ListCmd) SetFlags(f *flag.FlagSet) {
	// f.IntVar(&p.maxLineShift, "lines", 50, "allowed number of lines for a message can shift")
}

func (p *ListCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() < 1 {
		log.Fatal("Expects 1 or more checkfile arguments")
	}

	for _, fn := range fs.Args() {
		f, err := os.Open(fn)
		if err != nil {
			log.Fatalf("Failed to read '%s' - %s", fn, err)
		}
		defer f.Close()

		chk, err := checkstyle.Decode(f)
		if err != nil {
			log.Fatalf("Failed to parse '%s' - %s", fn, err)
		}

		for _, f := range chk.File {
			for _, e := range f.Error {
				fmt.Printf("%s on %s:%d - %s\n", e.Severity, f.Name, e.Line, e.Message)
			}
		}
	}

	return subcommands.ExitSuccess
}
