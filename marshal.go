package sxencode

import (
	"bytes"
)

func Marshal(v any) ([]byte, error) {
	var data bytes.Buffer
	enc := NewEncoder(&data)
	enc.Encode(v)
	return data.Bytes(), nil
}
