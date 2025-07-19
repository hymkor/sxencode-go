sxencode - Golang Encoder to S-Expression
=========================================

`sxencode` is a Go library that encodes arbitrary Go values into **Lisp-style S-expressions**.  
It is useful for interfacing with Lisp systems, structured debugging, or human-readable serialization.  
The output is readable by standard Lisp implementations including Common Lisp and ISLisp.


Example
-------

The following is a minimal Go program that encodes a Go struct and the encoder itself into S-expressions:

```example.go
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
```


## Output of `go run example.go`

The output of the above program is a pair of S-expressions representing the encoded data:

```make example|
go run example.go
((struct Foo)(Bar "hogehoge")(Baz 0.1)(Qux #(1 2 3 4))(Quux (("ahaha" 1)("ihihi" 2)("ufufu" 3)))(Quuux "a\"\\
    b"))
((struct Encoder)(TypeNotFound ""))
```

### Output Format

* A struct is encoded as a list beginning with `(struct <TypeName>)`, followed by pairs of field name and value:
  Example: `((struct Foo)(Bar "hogehoge") ...)`.
  The field names and type name are emitted as symbols, making them easy to extract with `(assoc)` in Lisp.

* A map is encoded as a list of `(key value)` pairs.
  If the keys are strings, note that in many Lisp dialects, `(assoc)` with a string key won't match unless `equal` is used.
  For this reason, helper functions like `field` (shown later) may be necessary.

* A slice or array is encoded using vector notation: `#(1 2 3 4)`.

These conventions make it easy to parse the output in various Lisp systems while retaining structural information.

## Reading the output in SBCL

The following shows how the output can be parsed and validated using [SBCL (Steel Bank Common Lisp)][SBCL].

[SBCL]: https://www.sbcl.org/ 

```make test-sbcl |
go run example.go | sbcl --load "test-sbcl.lsp"
This is SBCL 2.5.6, an implementation of ANSI Common Lisp.
More information about SBCL is available at <http://www.sbcl.org/>.

SBCL is free software, provided as is, with absolutely no warranty.
It is mostly in the public domain; some portions are provided under
BSD-style licenses.  See the CREDITS and COPYING files in the
distribution for more information.
PASS: (test (FIELD 'STRUCT DATA) FOO)
PASS: (test (FIELD 'BAR DATA) "hogehoge")
PASS: (test (FIELD 'BAZ DATA) 0.1)
PASS: (test (FIELD 'QUX DATA) #(1 2 3 4))
PASS: (test (FIELD "ahaha" M) 1)
PASS: (test (FIELD "ihihi" M) 2)
PASS: (test (FIELD "ufufu" M) 3)
PASS: (test (FIELD 'QUUUX DATA) "a\"\\
    b")
PASS: (test (FIELD 'STRUCT DATA) ENCODER)
PASS: (test (FIELD 'TYPENOTFOUND DATA) "")
* 
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
  (test (field 'struct data) 'Foo)
  (test (field 'bar data) "hogehoge")
  (test (field 'baz data) 0.1)
  (test (field 'qux data) #(1 2 3 4))
  (let ((m (field 'quux data)))
    (test (field "ahaha" m) 1)
    (test (field "ihihi" m) 2)
    (test (field "ufufu" m) 3))
  (test (field 'quuux data) "a\"\\
    b"))

(let ((data (read (standard-input) nil nil)))
  (test (field 'struct data) 'Encoder)
  (test (field 'typenotfound data) ""))
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
ISLisp>PASS: (test (FIELD (QUOTE STRUCT) DATA) FOO)
PASS: (test (FIELD (QUOTE BAR) DATA) "hogehoge")
PASS: (test (FIELD (QUOTE BAZ) DATA) 0.1)
PASS: (test (FIELD (QUOTE QUX) DATA) #(1 2 3 4))
PASS: (test (FIELD "ahaha" M) 1)
PASS: (test (FIELD "ihihi" M) 2)
PASS: (test (FIELD "ufufu" M) 3)
PASS: (test (FIELD (QUOTE QUUUX) DATA) "a\"\\
    b")
PASS: (test (FIELD (QUOTE STRUCT) DATA) ENCODER)
PASS: (test (FIELD (QUOTE TYPENOTFOUND) DATA) "")
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
PASS: (test (FIELD 'STRUCT DATA) FOO)
PASS: (test (FIELD 'BAR DATA) "hogehoge")
PASS: (test (FIELD 'BAZ DATA) 0.100000)
PASS: (test (FIELD 'QUX DATA) #(1 2 3 4))
PASS: (test (FIELD "ahaha" M) 1)
PASS: (test (FIELD "ihihi" M) 2)
PASS: (test (FIELD "ufufu" M) 3)
PASS: (test (FIELD 'QUUUX DATA) "a\"\\\n\tb")
PASS: (test (FIELD 'STRUCT DATA) ENCODER)
PASS: (test (FIELD 'TYPENOTFOUND DATA) "")
```

## Summary

* `sxencode` converts Go values into Lisp-compatible S-expressions
* Output is compatible with multiple Lisp systems: Common Lisp, ISLisp (OKI, gmnlisp)
* Includes reusable test logic in `test.lsp`
* Ideal for inter-language data exchange and human-readable serialization

Pull requests and feedback are welcome!
