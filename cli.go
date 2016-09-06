package cli

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

var ops = make(map[string]func())

func Define(op string, fn func()) {
	ops[op] = fn
}

func Usage(args ...interface{}) {
	if len(args) > 0 {
		fmt.Fprintln(os.Stderr, args...)
	}

	var kk []string
	for k := range ops {
		kk = append(kk, k)
	}
	sort.Strings(kk)
	fmt.Fprintf(os.Stderr, "ops: %v\n", kk)
}

func ParseFlag(onFlagParsed ...func() bool) {
	flag.Parse()

	for _, fn := range onFlagParsed {
		if !fn() {
			flag.Usage()
			os.Exit(2)
		}
	}
}

func findOp(op string) (fn func(), ok bool) {
	if v, ok := ops[op]; ok {
		return v, true
	}

	for k, v := range ops {
		if !strings.HasPrefix(k, op) {
			continue
		}
		if fn == nil {
			fn = v
			ok = true
		} else {
			fmt.Fprintf(os.Stderr, "matched multiple ops")
			os.Exit(2)
		}
	}
	return fn, ok
}

func Main() {
	if len(os.Args) < 2 || os.Args[1] == "-h" {
		Usage()
		os.Exit(2)
	}

	fn, ok := findOp(os.Args[1])
	if !ok {
		Usage()
		os.Exit(2)
	}

	os.Args = os.Args[1:]
	fn()
}
