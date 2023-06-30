package main

import (
	"flag"
	"os"

	"github.com/pascaldekloe/pile"
)

var (
	reverse          bool
	terminator       string
	requires, denies []string
)

func main() {
	flag.BoolFunc("require", "The output is suppressed with exit status 4 when the `string` is not present amongst the input.", func(s string) error {
		requires = append(requires, s)
		return nil
	})
	flag.BoolFunc("deny", "The output is suppressed with exit status 4 when the `string` is present amongst the input.", func(s string) error {
		denies = append(denies, s)
		return nil
	})
	flag.StringVar(&terminator, "terminate", "\n", "End each output with the `string`.")
	flag.BoolVar(&reverse, "reverse", false, "Print output in descending order.")
	flag.Parse()

	var args pile.Set[string]

	// read
	for _, s := range flag.Args() {
		args.Insert(s)
	}

	// filter
	for _, s := range requires {
		if !args.Find(s) {
			os.Exit(4)
		}
	}
	for _, s := range denies {
		if args.Find(s) {
			os.Exit(4)
		}
	}

	// print
	if reverse {
		for c, ok := args.Least(); ok; ok = c.Descend() {
			os.Stdout.WriteString(c.Key())
			os.Stdout.WriteString(terminator)
		}
	} else {
		for c, ok := args.Least(); ok; ok = c.Ascend() {
			os.Stdout.WriteString(c.Key())
			os.Stdout.WriteString(terminator)
		}
	}
}
