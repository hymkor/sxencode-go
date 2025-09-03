* Refactored the `parser` subpackage:
  * Removed the callback functions `Int`, `BigInt`, and `Float`, and added a unified `Number` callback that takes a string argument. As a result, importing `math/big` is no longer required by callers.
  * Changed the `String` callback specification to receive the raw string before backslash interpretation or newline conversion. This allows callers to handle backslash sequences in their own way.

v0.2.3
------
Aug 31, 2025

- Fixed an issue in the parser where a standalone backslash (`\`) inside a string literal (neither `\\` nor `\"`) was being dropped. Although such a literal is invalid in ISLisp, many other Lisp systems accept it, and dropping the backslash made it impossible to handle properly.

v0.2.2
------
Aug 22, 2025

- Updated `NewDecoder(r)` so that if `r` does not implement `io.RuneScanner`, it is automatically wrapped with `bufio.NewReader(r)`. As a result, users no longer need to wrap it manually.
- Updated the parser to match gmnlisp v0.7.22. Added fields such as `Function`, `Unquote`, `Quote`, and `Quasi` to the `parser.Parser` struct, allowing customization of how expressions like `#'`, `'`, `` ` ``, and `,` are turned into S-expressions (this does not affect the `sxencode` package specification itself).

v0.2.1
------
Jul 24, 2025

- Symbol can now be decoded into a string variable.

v0.2.0
------
Jul 23, 2025

- Made the delimiters for slices configurable via VectorOpen and VectorClose. The default was changed from `#(` and `)` to `(` and `)`.
- Implemented the `noname` tag option, which serializes the field as a standalone value instead of a `(name value)` pair.
- Removed the `(struct NAME)` format for struct names, as it was ultimately unused during decoding.
- Removed the `Name` type.
- Struct fields can now be decoded from both `(SYMBOL value)` and `("STRING" value)` forms.

v0.1.0
-------
Jul 23, 2025

- Implement `Decoder` and `Unmarshal`

v0.0.3
------
Jul 21, 2025

- Added support for struct tags `sxpr:"NAME,omitempty"` and `sxpr:",omitempty"` to omit fields with zero values from the S-expression output.
- Added support for `sxpr:"-"` struct tags to exclude fields from S-expression output entirely.
- Added support for a field of type ``sxencode.Name `sxpr:"SYMBOL"` `` in a struct to specify the symbol used in the struct header `(struct SYMBOL)`, similar to how "encoding/xml" works.
- Changed slice output to use the Lisp vector literal syntax `#(....)`.
- Removed the `ArrayHeader` and `ArrayIndex` fields from `Encoder`.
- Added support for calling the function set in `OnTypeNotSupported` when a type is not supported.
- For maps and structs, keys or field names are now omitted when their corresponding S-expression values are absent.
- For slices, if an element's S-expression is absent, `nil` is emitted instead.
- Added support for the struct tag `sxpr:"NAME"` to override the field name in S-expression output

v0.0.2
------

- Only `"` and `\` are escaped in string literals; other control characters (such as \n, \t, \r, etc.) are now output as raw characters
- Changed the package URL from `github.com/hymkor/sxencode` to `github.com/hymkor/sxencode-go`
- Changed struct type notation from `(struct-name NAME)` to `(struct NAME)`
- Modified `(*Encoder) Encode` to use the result of the `Sexpression() string` method, if implemented by an element in the input data
- Implemented the `Marshal` function

v0.0.1
------

- First release
