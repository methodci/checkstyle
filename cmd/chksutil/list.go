package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/google/subcommands"
	"github.com/methodci/checkstyle"
)

type listCmd struct {
	levelFilter string
}

func (*listCmd) Name() string     { return "list" }
func (*listCmd) Synopsis() string { return "list contents of one or more checkstyle files" }
func (*listCmd) Usage() string {
	return `list [<file>...]:
	list contents of one or more checkstyle files.
`
}

func (p *listCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.levelFilter, "severity", "", "comma separated list of severity levels to list - displays all on empty")
}

func (p *listCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() < 1 {
		log.Println("Expects 1 or more checkfile arguments")
		return subcommands.ExitUsageError
	}

	levels := strings.Split(p.levelFilter, ",")

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
				if (len(levels) == 1 && levels[0] == "") || contains(levels, string(e.Severity)) {
					fsev := formatSeverity(e.Severity)
					fmt.Printf("%s on %s:%d - %s\n", fsev("%s", e.Severity), f.Name, e.Line, e.Message)
				}
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

func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}
