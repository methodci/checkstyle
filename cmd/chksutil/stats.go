package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/google/subcommands"
	"github.com/methodci/checkstyle"
)

type statsCmd struct {
	top int
}

func (*statsCmd) Name() string     { return "stats" }
func (*statsCmd) Synopsis() string { return "get stats about one or more checkstyle file" }
func (*statsCmd) Usage() string {
	return `stats <file>:
	list contents of one or more checkstyle files.
`
}
func (p *statsCmd) SetFlags(f *flag.FlagSet) {
	f.IntVar(&p.top, "top", -1, "show only the top n files per severity, -1 means all")
}

func (p *statsCmd) Execute(_ context.Context, fs *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	if fs.NArg() < 1 {
		log.Println("Expects 1 or more checkfile arguments")
		return subcommands.ExitUsageError
	}

	typeCount := make(map[checkstyle.SeverityLevel]*struct {
		count int
		files map[string]int
	})

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
				if o, ok := typeCount[e.Severity]; !ok {
					typeCount[e.Severity] = &struct {
						count int
						files map[string]int
					}{
						count: 1,
						files: map[string]int{f.Name: 1},
					}
				} else {
					o.count++
					o.files[f.Name]++
				}
			}
		}
	}

	s := mapToSlice(typeCount)
	sort.Slice(s, func(i, j int) bool {
		return s[i].val.count > s[j].val.count
	})

	for _, v := range s {

		fmt.Printf("%s: %d total\n", formatSeverity(v.key)(string(v.key)), v.val.count)

		fs := mapToSlice(v.val.files)
		sort.Slice(fs, func(i, j int) bool {
			return fs[i].val > fs[j].val
		})

		counted := 0
		for c, f := range fs {
			if p.top > -1 && c >= p.top {
				break
			} else {
				fmt.Printf("\t%s : %d\n", f.key, f.val)
				counted += f.val
			}
		}

		if counted < v.val.count && p.top > 0 {
			fmt.Printf("\t... +%d more in %d files\n", v.val.count-counted, len(v.val.files)-p.top)
		}
	}

	return subcommands.ExitSuccess
}

func mapToSlice[T comparable, U any](s map[T]U) []struct {
	key T
	val U
} {
	var r []struct {
		key T
		val U
	}
	for k, v := range s {
		r = append(r, struct {
			key T
			val U
		}{k, v})
	}
	return r
}
