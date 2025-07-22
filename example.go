//go:build run

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/hymkor/sxencode-go"
)

var flagWarn = flag.Bool("w", false, "warning")

func mains() error {
	type Foo struct {
		Bar   string         `sxpr:"bar"`
		Baz   float64        `sxpr:"baz"`
		Qux   []int          `sxpr:"qux"`
		Quux  map[string]int `sxpr:"quux"`
		Quuux string         `sxpr:"quuux"`
		Corge func()         `sxpr:"corge"`
	}

	value := &Foo{
		Bar:   "hogehoge",
		Baz:   0.1,
		Qux:   []int{1, 2, 3, 4},
		Quux:  map[string]int{"ahaha": 1, "ihihi": 2, "ufufu": 3},
		Quuux: "a\"\\\n\tb",
		Corge: nil,
	}

	sxencode.VectorOpen = "#(" // for CommonLisp and ISLisp

	sxpr, err := sxencode.Marshal(value)
	if err != nil {
		return err
	}

	fmt.Println(string(sxpr))

	var clone Foo
	err = sxencode.Unmarshal(sxpr, &clone)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(value, &clone) {
		return errors.New("encode or decode failed")
	}
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
