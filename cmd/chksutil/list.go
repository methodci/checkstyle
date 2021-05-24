package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/google/subcommands"
	"github.com/methodci/checkstyle"
)

type listCmd struct {
	maxLineShift int
}

func (*listCmd) Name() string     { return "list" }
func (*listCmd) Synopsis() string { return "list contents of one or more checkstyle files" }
func (*listCmd) Usage() string {
	return `list <left-file> <right-file>:
	list contents of one or more checkstyle files.
`
}

func (p *listCmd) SetFlags(f *flag.FlagSet) {
	// f.IntVar(&p.maxLineShift, "lines", 50, "allowed number of lines for a message can shift")
}

func (p *listCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() < 1 {
		log.Println("Expects 1 or more checkfile arguments")
		return subcommands.ExitUsageError
	}

	for _, fn := range fs.Args() {
		f, err := os.Open(fn)
		if err != nil {
			log.Printf("Failed to read '%s' - %s", fn, err)
			return subcommands.ExitFailure
		}
		defer f.Close()

		chk, err := checkstyle.Decode(f)
		if err != nil {
			log.Printf("Failed to parse '%s' - %s", fn, err)
			return subcommands.ExitFailure
		}

		for _, f := range chk.File {
			for _, e := range f.Error {
				fsev := formatSeverity(e.Severity)
				fmt.Printf("%s on %s:%d - %s\n", fsev("%s", e.Severity), f.Name, e.Line, e.Message)
			}
		}
	}

	return subcommands.ExitSuccess
}

func formatSeverity(s checkstyle.SeverityLevel) func(string, ...interface{}) string {
	switch s {
	case checkstyle.SeverityError:
		return color.RedString
	case checkstyle.SeverityWarning:
		return color.YellowString
	case checkstyle.SeverityInfo:
		return color.CyanString
	}

	return fmt.Sprintf
}
