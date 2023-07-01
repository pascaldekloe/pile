package main

import (
	"flag"
	"os"

	"github.com/pascaldekloe/pile"
)

var (
	dedupe           bool
	reverse          bool
	terminator       string
	requires, denies []string
)

func main() {
	flag.Func("require", "Output is suppressed with exit status 4 when the `string` equals\nno input argument. Multiple -require options can be used in\nconjunction.", func(s string) error {
		requires = append(requires, s)
		return nil
	})
	flag.Func("deny", "Output is suppressed with exit status 3 when the `string` equals\nan input argument. Multiple -deny options can be used in\nconjunction.", func(s string) error {
		denies = append(denies, s)
		return nil
	})
	flag.StringVar(&terminator, "terminate", "\n", "Print `string` after each output.")
	flag.BoolVar(&dedupe, "dedupe", false, "Print duplicated input only once like uniq(1) does.")
	flag.BoolVar(&reverse, "reverse", false, "Print output in descending order.")
	flag.Parse()

	var args pile.Map[string, int]

	// read
	for _, s := range flag.Args() {
		n, _ := args.Find(s)
		args.Put(s, n+1)
	}

	// filter
	for _, s := range requires {
		if _, ok := args.Find(s); !ok {
			os.Exit(4)
		}
	}
	for _, s := range denies {
		if _, ok := args.Find(s); !ok {
			os.Exit(3)
		}
	}

	// print
	if reverse {
		for c, ok := args.Least(); ok; ok = c.Descend() {
			print(c)
		}
	} else {
		for c, ok := args.Least(); ok; ok = c.Ascend() {
			print(c)
		}
	}
}

func print(c pile.Cursor[string, int]) {
	n := 1
	if !dedupe {
		n = c.Value()
	}
	for ; n > 0; n-- {
		os.Stdout.WriteString(c.Key())
		os.Stdout.WriteString(terminator)
	}
}
