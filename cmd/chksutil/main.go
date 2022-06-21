package main

import (
	"context"
	"flag"
	"os"

	"github.com/fatih/color"
	"github.com/google/subcommands"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(&diffCmd{}, "")
	subcommands.Register(&listCmd{}, "")
	subcommands.Register(&statsCmd{}, "")

	nocolor := flag.Bool("no-color", false, "Disable colorized output")
	flag.Parse()
	if *nocolor {
		color.NoColor = true
	}

	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
