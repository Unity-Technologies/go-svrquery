package common

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

// WireEncoder is an interface which allows for different query implementations
// to write data to a byte buffer in a specific format.
type WireEncoder interface {
	WriteString(resp *bytes.Buffer, s string) error
	Write(resp *bytes.Buffer, v interface{}) error
}

// Encoder is a struct which implements proto.WireEncoder
type Encoder struct{}

// WriteString writes a string to the provided buffer.
func (e *Encoder) WriteString(resp *bytes.Buffer, s string) error {
	if err := binary.Write(resp, binary.BigEndian, byte(len(s))); err != nil {
		return err
	}

	return binary.Write(resp, binary.BigEndian, []byte(s))
}

// Write writes arbitrary data to the provided buffer.
func (e *Encoder) Write(resp *bytes.Buffer, v interface{}) error {
	return binary.Write(resp, binary.BigEndian, v)
}

// WireWrite writes the provided data to resp with the provided WireEncoder w.
func WireWrite(resp *bytes.Buffer, w WireEncoder, data interface{}) error {
	t := reflect.TypeOf(data)
	vs := reflect.Indirect(reflect.ValueOf(data))
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		v := vs.FieldByName(f.Name)

		// Dereference pointer
		if f.Type.Kind() == reflect.Ptr {
			if v.IsNil() {
				continue
			}
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.Struct:
			if err := WireWrite(resp, w, v.Interface()); err != nil {
				return err
			}

		case reflect.String:
			if err := w.WriteString(resp, v.String()); err != nil {
				return err
			}

		default:
			if err := w.Write(resp, v.Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}
