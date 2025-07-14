package sxencode

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type Encoder struct {
	w            io.Writer
	ArrayHeader  string
	ArrayIndex   bool
	TypeNotFound string
}

func (enc *Encoder) writeByte(b byte) error {
	_, err := enc.w.Write([]byte{b})
	return err
}

func (enc *Encoder) write(b []byte) (int, error) {
	return enc.w.Write(b)
}

func (enc *Encoder) writeString(s string) (int, error) {
	return io.WriteString(enc.w, s)
}

type Sexpressioner interface {
	Sexpression() string
}

var toLispString = strings.NewReplacer(
	`"`, `\"`,
	`\`, `\\`,
)

func (enc *Encoder) encode(value reflect.Value) {
	k := value.Kind()
	if value.CanInterface() {
		if v, ok := value.Interface().(Sexpressioner); ok {
			io.WriteString(enc.w, v.Sexpression())
			return
		}
	}
	switch k {
	case reflect.Interface, reflect.Pointer:
		enc.encode(value.Elem())
	case reflect.Struct:
		types := value.Type()
		enc.writeByte('(')
		if name := types.Name(); name != "" {
			fmt.Fprintf(enc.w, "(struct %s)", name)
		}
		fields := reflect.VisibleFields(types)
		for i, t := range fields {
			if t.IsExported() {
				fmt.Fprintf(enc.w, "(%s ", t.Name)
				enc.encode(value.Field(i))
				enc.writeByte(')')
			}
		}
		enc.writeByte(')')
	case reflect.String:
		enc.writeByte('"')
		io.WriteString(enc.w, toLispString.Replace(value.String()))
		enc.writeByte('"')
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprint(enc.w, value.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fmt.Fprint(enc.w, value.Uint())
	case reflect.Float32, reflect.Float64:
		fmt.Fprint(enc.w, value.Float())
	case reflect.Array, reflect.Slice:
		if enc.ArrayHeader != "" {
			fmt.Fprintf(enc.w, "(%s ", enc.ArrayHeader)
		} else {
			enc.writeByte('(')
		}
		if n := value.Len(); n > 0 {
			i := 0
			for {
				if enc.ArrayIndex {
					fmt.Fprintf(enc.w, "(%d ", i)
					enc.encode(value.Index(i))
					enc.writeByte(')')
				} else {
					enc.encode(value.Index(i))
				}
				if i++; i >= n {
					break
				}
				enc.writeByte(' ')
			}
		}

		enc.writeByte(')')
	case reflect.Map:
		iter := value.MapRange()
		enc.writeByte('(')
		for iter.Next() {
			enc.writeByte('(')
			enc.encode(iter.Key())
			enc.writeByte(' ')
			enc.encode(iter.Value())
			enc.writeByte(')')
		}
		enc.writeByte(')')
	case reflect.Bool:
		if value.Bool() {
			enc.writeByte('t')
		} else {
			enc.writeString("nil")
		}
	default:
		if enc.TypeNotFound != "" {
			fmt.Fprintf(enc.w, "(%s %#v)", enc.TypeNotFound, value.String())
		}
	}
}

func (enc *Encoder) Encode(v any) {
	enc.encode(reflect.ValueOf(v))
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}
