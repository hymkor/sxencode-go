sxencode - Golang Encoder to S-Expression
=========================================

`sxencode` is a Go library that encodes arbitrary Go values into **Lisp-style S-expressions**.  
It is useful for interfacing with Lisp systems, structured debugging, or human-readable serialization.  
The output is readable by standard Lisp implementations including Common Lisp and ISLisp.

```go doc -all|
package sxencode // import "github.com/hymkor/sxencode-go"


VARIABLES

var (
    // Delimiters for vector (array/slice) literals in S-expression encoding.
    VectorOpen  = "("
    VectorClose = ")"
)

FUNCTIONS

func Marshal(v any) ([]byte, error)
func Unmarshal(data []byte, v any) error

TYPES

type Decoder struct {
    OnTypeNotSupported func(any, reflect.Value) error
    // Has unexported fields.
}

func NewDecoder(r io.RuneScanner) *Decoder

func (D *Decoder) Decode(v any) error

type Encoder struct {
    OnTypeNotSupported func(reflect.Value) (string, error)
    // Has unexported fields.
}

func NewEncoder(w io.Writer) *Encoder

func (enc *Encoder) Encode(v any) error

type Sexpressioner interface {
    Sexpression() string
}

```

Example
-------

The following is a minimal Go program that encodes a Go struct and the encoder itself into S-expressions:

```example.go
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
```


## Output of `go run example.go`

The output of the above program is a pair of S-expressions representing the encoded data:

```make example|
go run example.go
((bar "hogehoge")(baz 0.1)(qux (1 2 3 4))(quux (("ahaha" 1)("ihihi" 2)("ufufu" 3)))(quuux "a\"\\
    b"))
```

### Output Format

* A struct is encoded as a list of field name and value pairs:
  Example: `((Bar "hogehoge") (Baz 42) ...)`.
  The field names are emitted as symbols, making them easy to extract with `(assoc)` in Lisp.
  The type name is no longer included in the output.

* A map is encoded as a list of `(key value)` pairs.
  If the keys are strings, note that in many Lisp dialects, `(assoc)` with a string key won't match unless `equal` is used.
  For this reason, helper functions like `field` (shown later) may be necessary.

* A slice or array is encoded using vector notation: `#(1 2 3 4)`.

These conventions make it easy to parse the output in various Lisp systems while retaining structural information.

## Reading the output in SBCL

The following shows how the output can be parsed and validated using [SBCL (Steel Bank Common Lisp)][SBCL].

[SBCL]: https://www.sbcl.org/ 

```make test-sbcl |
go run example.go | sbcl --script "test-sbcl.lsp"
PASS: (test (FIELD 'BAR DATA) "hogehoge")
PASS: (test (FIELD 'BAZ DATA) 0.1)
PASS: (test (FIELD 'QUX DATA) (1 2 3 4))
PASS: (test (FIELD "ahaha" M) 1)
PASS: (test (FIELD "ihihi" M) 2)
PASS: (test (FIELD "ufufu" M) 3)
PASS: (test (FIELD 'QUUUX DATA) "a\"\\
    b")
```

### Lisp files

These are the supporting Lisp files used for the SBCL test:

#### `test-sbcl.lsp`

```test-sbcl.lsp
;; Define compatibility functions to match ISLisp's standard input/output access
(defun standard-input () *standard-input*)
(defun standard-output () t)
(load "test.lsp")
```

#### `test.lsp`

```test.lsp
(defmacro test (source expect)
  (let ((result (gensym)))
    `(let ((,result ,source))
       (if (equalp ,result ,expect)
           (format (standard-output) "PASS: (test ~S ~S)~%"
                   (quote ,source)
                   ,expect)
           (format (standard-output) "FAIL: (test ~S ~S)~%  but ~S~%"
                   (quote ,source)
                   ,expect
                   ,result)
       ))))

(defun field (key m)
  (and
    m
    (consp m)
    (if (equal (car (car m)) key)
      (car (cdr (car m)))
      (field key (cdr m)))))

(let ((data (read (standard-input) nil nil)))
  (test (field 'bar data) "hogehoge")
  (test (field 'baz data) 0.1)
  (test (field 'qux data) '(1 2 3 4))
  (let ((m (field 'quux data)))
    (test (field "ahaha" m) 1)
    (test (field "ihihi" m) 2)
    (test (field "ufufu" m) 3))
  (test (field 'quuux data) "a\"\\
    b"))
```


## Reading the output in OKI ISLisp

The same logic can be used in [OKI ISLisp](https://openlab.jp/islisp/), by loading the wrapper file `test-oki.lsp`:

```sh
go run example.go > sample.log
echo "(load \"test-oki.lsp\")" | islisp
```

Execution result:

```make test-oki|
go run example.go > sample.log
echo "(load \"test-oki.lsp\")" | islisp
> ISLisp  Version 0.80 (1999/02/25)
>
ISLisp>PASS: (test (FIELD (QUOTE BAR) DATA) "hogehoge")
PASS: (test (FIELD (QUOTE BAZ) DATA) 0.1)
PASS: (test (FIELD (QUOTE QUX) DATA) (1 2 3 4))
PASS: (test (FIELD "ahaha" M) 1)
PASS: (test (FIELD "ihihi" M) 2)
PASS: (test (FIELD "ufufu" M) 3)
PASS: (test (FIELD (QUOTE QUUUX) DATA) "a\"\\
    b")
T
ISLisp>
```

#### `test-oki.lsp`

```test-oki.lsp
(defun equalp (x y)
  (equal x y))

(with-open-input-file
  (fd "sample.log")
  (with-standard-input
    fd
    (load "test.lsp")))
```


## Also works with gmnlisp

The same test.lsp file can be run directly using [gmnlisp] — no modifications are needed, since [gmnlisp] is designed to be largely compatible with ISLisp and supports the same input/output conventions.

[gmnlisp]: https://github.com/hymkor/gmnlisp

```sh
go run example.go | gmnlisp test.lsp
```

Result:

```make test-gmnlisp|
go run example.go | gmnlisp test.lsp
PASS: (test (FIELD 'BAR DATA) "hogehoge")
PASS: (test (FIELD 'BAZ DATA) 0.100000)
PASS: (test (FIELD 'QUX DATA) (1 2 3 4))
PASS: (test (FIELD "ahaha" M) 1)
PASS: (test (FIELD "ihihi" M) 2)
PASS: (test (FIELD "ufufu" M) 3)
PASS: (test (FIELD 'QUUUX DATA) "a\"\\\n\tb")
```

## Summary

* `sxencode` converts Go values into Lisp-compatible S-expressions
* Output is compatible with multiple Lisp systems: Common Lisp, ISLisp (OKI, gmnlisp)
* Includes reusable test logic in `test.lsp`
* Ideal for inter-language data exchange and human-readable serialization

Pull requests and feedback are welcome!
