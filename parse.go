package sxencode

import (
	"math/big"

	"github.com/hymkor/sxencode-go/parser"
)

type Symbol struct {
	Value string
}

type Cons struct {
	Car any
	Cdr any
}

var parser1 = &parser.Parser[any]{
	Cons:   func(car, cdr any) any { return &Cons{Car: car, Cdr: cdr} },
	Int:    func(n int64) any { return n },
	BigInt: func(n *big.Int) any { return n },
	Float:  func(f float64) any { return f },
	String: func(s string) any { return s },
	Array: func(list []any, dim []int) any {
		array := make([]any, len(list))
		for i, v := range list {
			array[i] = v
		}
		return array
	},
	Keyword: func(s string) any { return s },
	Rune:    func(r rune) any { return r },
	Symbol:  func(s string) any { return Symbol{Value: s} },
	Null:    func() any { return nil },
	True:    func() any { return true },
}
