package main

import (
	"context"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/google/subcommands"
	"github.com/methodci/checkstyle"
)

const (
	chkFixed   = "fixed"
	chkCreated = "created"
)

type diffCmd struct {
	maxLineShift int
	checkstyle   string
}

func (*diffCmd) Name() string     { return "diff" }
func (*diffCmd) Synopsis() string { return "diff two checkstyle files" }
func (*diffCmd) Usage() string {
	return `diff <left-file> <right-file>:
	diff two checkstyle files.
`
}

func (p *diffCmd) SetFlags(f *flag.FlagSet) {
	f.IntVar(&p.maxLineShift, "lines", 50, "allowed number of lines for a message can shift")
	f.StringVar(&p.checkstyle, "output-checkstyle", "", "output as checkstyle - options: "+chkFixed+" "+chkCreated)
}

func (p *diffCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() != 2 {
		log.Println("Expects exactly 2 checkfile file arguments")
		return subcommands.ExitUsageError
	}

	if p.checkstyle != "" && p.checkstyle != chkFixed && p.checkstyle != chkCreated {
		log.Println("checkstyle option expects value in  " + chkFixed + " " + chkCreated)
		return subcommands.ExitUsageError
	}

	fn1 := fs.Arg(0)
	fn2 := fs.Arg(1)

	f, err := os.Open(fn1)
	if err != nil {
		log.Printf("Failed to read '%s' - %s", fn1, err)
		return subcommands.ExitFailure
	}

	chk1, err := checkstyle.Decode(f)
	if err != nil {
		log.Printf("Failed to parse '%s' - %s", fn1, err)
		return subcommands.ExitFailure
	}

	f2, err := os.Open(fn2)
	if err != nil {
		log.Printf("Failed to read '%s' - %s", fn2, err)
		return subcommands.ExitFailure
	}

	chk2, err := checkstyle.Decode(f2)
	if err != nil {
		log.Printf("Failed to parse '%s' - %s", fn2, err)
		return subcommands.ExitFailure
	}

	fixedErr, newErr := checkstyle.Diff(chk1, chk2, checkstyle.DiffOptions{MaxLineDiff: p.maxLineShift})

	if p.checkstyle != "" {
		enc := xml.NewEncoder(os.Stdout)
		enc.Indent("", "\t")
		switch p.checkstyle {
		case chkFixed:
			enc.Encode(fixedErr)
		case chkCreated:
			enc.Encode(newErr)
		}

		os.Stdout.WriteString("\n")

		return subcommands.ExitSuccess
	}

	for _, f := range fixedErr.File {
		for _, e := range f.Error {
			fmt.Printf("%s %s on %s:%d - %s\n", color.GreenString("%s %s", "Fixed", e.Severity), color.MagentaString(e.Source), f.Name, e.Line, e.Message)
		}
	}

	for _, f := range newErr.File {
		for _, e := range f.Error {
			fsev := formatSeverity(e.Severity)
			fmt.Printf("%s %s on %s:%d - %s\n", fsev("%s %s", "Created", e.Severity), color.MagentaString(e.Source), f.Name, e.Line, e.Message)
		}
	}

	return subcommands.ExitSuccess
}
