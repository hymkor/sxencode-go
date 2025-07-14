//go:build run

package main

import (
	"fmt"
	"os"

	"github.com/hymkor/sxencode-go"
)

func main() {
	type Foo struct {
		Bar   string
		Baz   float64
		Qux   []int
		Quux  map[string]int
		Quuux string
	}

	value := &Foo{
		Bar:   "hogehoge",
		Baz:   0.1,
		Qux:   []int{1, 2, 3, 4},
		Quux:  map[string]int{"ahaha": 1, "ihihi": 2, "ufufu": 3},
		Quuux: "a\"\\\n\tb",
	}

	enc := sxencode.NewEncoder(os.Stdout)

	enc.Encode(value)
	fmt.Println()

	enc.Encode(enc)
	fmt.Println()
}
