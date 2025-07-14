- Only `"` and `\` are escaped in string literals; other control characters (such as \n, \t, \r, etc.) are now output as raw characters

### Jul 14, 2025

- Changed the package URL from `github.com/hymkor/sxencode` to `github.com/hymkor/sxencode-go`
- Changed struct type notation from `(struct-name NAME)` to `(struct NAME)`
- Modified `(*Encoder) Encode` to use the result of the `Sexpression() string` method, if implemented by an element in the input data
- Implemented the `Marshal` function

### Jul 13, 2025

- First release
