//go:build run

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hymkor/sxencode-go"
)

var flagWarn = flag.Bool("w", false, "warning")

func main() {
	type Foo struct {
		Bar   string
		Baz   float64
		Qux   []int
		Quux  map[string]int
		Quuux string
		Corge func()
	}

	value := &Foo{
		Bar:   "hogehoge",
		Baz:   0.1,
		Qux:   []int{1, 2, 3, 4},
		Quux:  map[string]int{"ahaha": 1, "ihihi": 2, "ufufu": 3},
		Quuux: "a\"\\\n\tb",
		Corge: func() {},
	}

	enc := sxencode.NewEncoder(os.Stdout)
	if *flagWarn {
		enc.OnTypeNotFound = func(v reflect.Value) (string, error) {
			return "'not-support-type", nil
		}
	}

	enc.Encode(value)
	fmt.Println()

	enc.Encode(enc)
	fmt.Println()
}
