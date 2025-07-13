//go:build run

package main

import (
	"fmt"
	"os"

	"github.com/hymkor/sxencode"
)

func main() {
	type Foo struct {
		Name  string
		Value float64
		Array []int
		Map   map[string]int
	}

	value := &Foo{
		Name:  "hogehoge",
		Value: 0.1,
		Array: []int{1, 2, 3, 4},
		Map:   map[string]int{"ahaha": 1, "ihihi": 2, "ufufu": 3},
	}

	enc := sxencode.NewEncoder(os.Stdout)
	// enc.ArrayHeader = "array"
	// enc.ArrayIndex = true
	// enc.TypeNotFound = "type-not-found"

	enc.Encode(value)
	fmt.Println()

	enc.Encode(enc)
	fmt.Println()
}
