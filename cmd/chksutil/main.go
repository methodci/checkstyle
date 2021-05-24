package main

import (
	"context"
	"flag"
	"os"

	"github.com/fatih/color"
	"github.com/google/subcommands"
)

func init() {
	flag.BoolVar(&color.NoColor, "no-color", false, "Disable colorized output")
}

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(&DiffCmd{}, "")
	subcommands.Register(&ListCmd{}, "")

	flag.Parse()
	ctx := context.Background()
	os.Exit(int(subcommands.Execute(ctx)))
}
