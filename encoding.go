package sxencode

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type Encoder struct {
	w              io.Writer
	OnTypeNotFound func(reflect.Value) (string, error)
}

func (enc *Encoder) writeByte(b byte) error {
	_, err := enc.w.Write([]byte{b})
	return err
}

func (enc *Encoder) write(b []byte) error {
	_, err := enc.w.Write(b)
	return err
}

func (enc *Encoder) writeString(s string) error {
	_, err := io.WriteString(enc.w, s)
	return err
}

type Sexpressioner interface {
	Sexpression() string
}

var toLispString = strings.NewReplacer(
	`"`, `\"`,
	`\`, `\\`,
)

func (enc *Encoder) encode(value reflect.Value) error {
	k := value.Kind()
	if value.CanInterface() {
		if v, ok := value.Interface().(Sexpressioner); ok {
			_, err := io.WriteString(enc.w, v.Sexpression())
			return err
		}
	}
	switch k {
	case reflect.Interface, reflect.Pointer:
		return enc.encode(value.Elem())
	case reflect.Struct:
		types := value.Type()
		if err := enc.writeByte('('); err != nil {
			return err
		}
		if name := types.Name(); name != "" {
			if _, err := fmt.Fprintf(enc.w, "(struct %s)", name); err != nil {
				return err
			}
		}
		fields := reflect.VisibleFields(types)
		for i, t := range fields {
			if t.IsExported() {
				if _, err := fmt.Fprintf(enc.w, "(%s ", t.Name); err != nil {
					return err
				}
				if err := enc.encode(value.Field(i)); err != nil {
					return err
				}
				if err := enc.writeByte(')'); err != nil {
					return err
				}
			}
		}
		return enc.writeByte(')')
	case reflect.String:
		if err := enc.writeByte('"'); err != nil {
			return err
		}
		if _, err := io.WriteString(enc.w, toLispString.Replace(value.String())); err != nil {
			return err
		}
		return enc.writeByte('"')
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, err := fmt.Fprint(enc.w, value.Int())
		return err
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		_, err := fmt.Fprint(enc.w, value.Uint())
		return err
	case reflect.Float32, reflect.Float64:
		_, err := fmt.Fprint(enc.w, value.Float())
		return err
	case reflect.Array, reflect.Slice:
		if err := enc.writeString("#("); err != nil {
			return err
		}
		if n := value.Len(); n > 0 {
			i := 0
			for {
				if err := enc.encode(value.Index(i)); err != nil {
					return err
				}
				if i++; i >= n {
					break
				}
				if err := enc.writeByte(' '); err != nil {
					return err
				}
			}
		}
		return enc.writeByte(')')
	case reflect.Map:
		iter := value.MapRange()
		enc.writeByte('(')
		for iter.Next() {
			if err := enc.writeByte('('); err != nil {
				return err
			}
			if err := enc.encode(iter.Key()); err != nil {
				return err
			}
			if err := enc.writeByte(' '); err != nil {
				return err
			}
			if err := enc.encode(iter.Value()); err != nil {
				return err
			}
			if err := enc.writeByte(')'); err != nil {
				return err
			}
		}
		return enc.writeByte(')')
	case reflect.Bool:
		if value.Bool() {
			return enc.writeByte('t')
		} else {
			return enc.writeString("nil")
		}
	default:
		if enc.OnTypeNotFound != nil {
			s, err := enc.OnTypeNotFound(value)
			if err != nil {
				return err
			}
			return enc.writeString(s)
		} else {
			return enc.writeString("nil")
		}
	}
	return nil
}

func (enc *Encoder) Encode(v any) error {
	return enc.encode(reflect.ValueOf(v))
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}
