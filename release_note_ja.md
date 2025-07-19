### Jul 20, 2025

- スライスの出力には、Lisp のベクタリテラル `#(....)` を用いるようにした

### Jul 15, 2025

- 文字列中の `"` と `\` 以外の制御文字（例：`\n`, `\t`, `\r`, `\b`, `\a`）はエスケープせず、生の文字として出力するようにした

### Jul 14, 2025

- パッケージの URL を `github.com/hymkor/sxencode` から `github.com/hymkor/sxencode-go` へ変更した
- `(struct-name NAME)` だった構造体の型名表記を `(struct NAME)` へ変更した
- `(*Encoder) Encode` に与えられたデータに含まれる要素に `Sexpression() string` というメソッドが実装されていた時、その要素のS式化にはそのメソッドの結果を利用するようにした。
- `Marshal` 関数を実装

### Jul 13, 2025

- 公開
