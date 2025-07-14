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
```


## Output of `go run example.go`

The output of the above program is a pair of S-expressions representing the encoded data:

```make example|
go run example.go
((struct Foo)(Name "hogehoge")(Value 0.1)(Array (1 2 3 4))(Map (("ahaha" 1)("ihihi" 2)("ufufu" 3))))
((struct Encoder)(ArrayHeader "")(ArrayIndex nil)(TypeNotFound ""))
```

### Output Format

* A struct is encoded as a list beginning with `(struct <TypeName>)`, followed by pairs of field name and value:
  Example: `((struct Foo)(Name "hogehoge") ...)`.
  The field names and type name are emitted as symbols, making them easy to extract with `(assoc)` in Lisp.

* A map is encoded as a list of `(key value)` pairs.
  If the keys are strings, note that in many Lisp dialects, `(assoc)` with a string key won't match unless `equal` is used.
  For this reason, helper functions like `field` (shown later) may be necessary.

* A slice or array is encoded as a plain list of elements: `(1 2 3 4)`.

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
PASS: (test (FIELD 'NAME DATA) "hogehoge")
PASS: (test (FIELD 'VALUE DATA) 0.1)
PASS: (test (FIELD 'ARRAY DATA) (1 2 3 4))
PASS: (test (FIELD "ahaha" M) 1)
PASS: (test (FIELD "ihihi" M) 2)
PASS: (test (FIELD "ufufu" M) 3)
PASS: (test (FIELD 'STRUCT DATA) ENCODER)
PASS: (test (FIELD 'ARRAYHEADER DATA) "")
PASS: (test (FIELD 'ARRAYINDEX DATA) NIL)
PASS: (test (FIELD 'TYPENOTFOUND DATA) "")
* 
```

### Lisp files

These are the supporting Lisp files used for the SBCL test:

#### `test-sbcl.lsp`

```lisp
(defun standard-input () *standard-input*)
(defun standard-output () t)
(load "test.lsp")
```

#### `test.lsp`

```lisp
(defmacro test (source expect)
  (let ((result (gensym)))
    `(let ((,result ,source))
       (if (equal ,result ,expect)
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
  (test (field 'struct-name data) 'Foo)
  (test (field 'name data) "hogehoge")
  (test (field 'value data) 0.1)
  (test (field 'array data) '(1 2 3 4))
  (let ((m (field 'map data)))
    (test (field "ahaha" m) 1)
    (test (field "ihihi" m) 2)
    (test (field "ufufu" m) 3)))

(let ((data (read (standard-input) nil nil)))
  (test (field 'struct-name data) 'Encoder)
  (test (field 'arrayheader data) "")
  (test (field 'arrayindex data) nil)
  (test (field 'typenotfound data) ""))
```


## Reading the output in OKI ISLisp

The same logic can be used in [OKI ISLisp](https://openlab.jp/islisp/), by loading the wrapper file `test-oki.lsp`:

```sh
go run example.go > sample.log
echo "(load \"test-oki.lsp\")" | islisp
```

Execution result:

```
> ISLisp  Version 0.80 (1999/02/25)
>
ISLisp>PASS: (test (FIELD (QUOTE STRUCT-NAME) DATA) FOO)
PASS: (test (FIELD (QUOTE NAME) DATA) "hogehoge")
PASS: (test (FIELD (QUOTE VALUE) DATA) 0.1)
PASS: (test (FIELD (QUOTE ARRAY) DATA) (1 2 3 4))
PASS: (test (FIELD "ahaha" M) 1)
PASS: (test (FIELD "ihihi" M) 2)
PASS: (test (FIELD "ufufu" M) 3)
PASS: (test (FIELD (QUOTE STRUCT-NAME) DATA) ENCODER)
PASS: (test (FIELD (QUOTE ARRAYHEADER) DATA) "")
PASS: (test (FIELD (QUOTE ARRAYINDEX) DATA) NIL)
PASS: (test (FIELD (QUOTE TYPENOTFOUND) DATA) "")
T
ISLisp>
```

#### `test-oki.lsp`

```lisp
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

```
PASS: (test (FIELD 'STRUCT-NAME DATA) FOO)
PASS: (test (FIELD 'NAME DATA) "hogehoge")
PASS: (test (FIELD 'VALUE DATA) 0.100000)
PASS: (test (FIELD 'ARRAY DATA) (1 2 3 4))
PASS: (test (FIELD "ahaha" M) 1)
PASS: (test (FIELD "ihihi" M) 2)
PASS: (test (FIELD "ufufu" M) 3)
PASS: (test (FIELD 'STRUCT-NAME DATA) ENCODER)
PASS: (test (FIELD 'ARRAYHEADER DATA) "")
PASS: (test (FIELD 'ARRAYINDEX DATA) NIL)
PASS: (test (FIELD 'TYPENOTFOUND DATA) "")
```

## Summary

* `sxencode` converts Go values into Lisp-compatible S-expressions
* Output is compatible with multiple Lisp systems: Common Lisp, ISLisp (OKI, gmnlisp)
* Includes reusable test logic in `test.lsp`
* Ideal for inter-language data exchange and human-readable serialization

Pull requests and feedback are welcome!
