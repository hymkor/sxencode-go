package sxencode

import (
	"reflect"
	"strings"
	"testing"
)

type foo struct {
	Bar   string
	Corge func()
}

func TestStruct(t *testing.T) {
	v := &foo{
		Bar:   "hoge",
		Corge: func() {},
	}
	b, err := Marshal(v)
	if err != nil {
		t.Fatal(err.Error())
	}
	expect := `((struct foo)(Bar "hoge"))`
	result := string(b)
	if expect != result {
		t.Fatalf("expect %v, but %v", expect, result)
	}

	var sbuf strings.Builder
	enc := NewEncoder(&sbuf)
	enc.OnTypeNotFound = func(v reflect.Value) (string, error) {
		return "not-support-type", nil
	}
	enc.Encode(v)
	result = sbuf.String()
	expect1 := `((struct foo)(Bar "hoge")(Corge not-support-type))`
	expect2 := `((struct foo)(Corge not-support-type)(Bar "hoge"))`
	if expect1 != result && expect2 != result {
		t.Fatalf("expect %v or %v, but %v", expect1, expect2, result)
	}
}

func TestMap(t *testing.T) {
	v := map[string]any{
		"bar": "hoge",
		"baz": func() {},
	}
	b, err := Marshal(v)
	if err != nil {
		t.Fatal(err.Error())
	}
	expect := `(("bar" "hoge"))`
	result := string(b)
	if expect != result {
		t.Fatalf("expect %v, but %v", expect, result)
	}

	var sbuf strings.Builder
	enc := NewEncoder(&sbuf)
	enc.OnTypeNotFound = func(v reflect.Value) (string, error) {
		return "not-support-type", nil
	}
	enc.Encode(v)
	result = sbuf.String()
	expect1 := `(("bar" "hoge")("baz" not-support-type))`
	expect2 := `(("baz" not-support-type)("bar" "hoge"))`
	if expect1 != result && expect2 != result {
		t.Fatalf("expect %v or %v, but %v", expect1, expect2, result)
	}
}
