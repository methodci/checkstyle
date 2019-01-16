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

type DiffCmd struct {
	maxLineShift int
}

func (*DiffCmd) Name() string     { return "diff" }
func (*DiffCmd) Synopsis() string { return "diff two checkstyle files" }
func (*DiffCmd) Usage() string {
	return `diff <left-file> <right-file>:
	list your notes.
  `
}

func (p *DiffCmd) SetFlags(f *flag.FlagSet) {
	f.IntVar(&p.maxLineShift, "lines", 50, "allowed number of lines for a message can shift")
}

func (p *DiffCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() != 2 {
		log.Fatal("Expects exactly 2 checkfile arguments")
	}

	f, err := os.Open(fs.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	chk1, err := checkstyle.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	f2, err := os.Open(fs.Arg(1))
	if err != nil {
		log.Fatal(err)
	}

	chk2, err := checkstyle.Decode(f2)
	if err != nil {
		log.Fatal(err)
	}

	fixedErr, newErr := checkstyle.Diff(chk1, chk2, checkstyle.DiffOptions{MaxLineDiff: 50})
	for _, f := range fixedErr.File {
		for _, e := range f.Error {
			fmt.Printf("Fixed %s on %s:%d - %s\n", e.Severity, f.Name, e.Line, e.Message)
		}
	}

	for _, f := range newErr.File {
		for _, e := range f.Error {
			fmt.Printf("Created %s on %s:%d - %s\n", e.Severity, f.Name, e.Line, e.Message)
		}
	}

	return subcommands.ExitSuccess
}
