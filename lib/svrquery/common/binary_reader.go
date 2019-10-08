package common

import (
	"bytes"
	"encoding/binary"
)

// BinaryReader is a reader which uses binary encoding.
type BinaryReader struct {
	buf   *bytes.Buffer
	order binary.ByteOrder
}

// NewBinaryReader returns a new BinaryReader which reads data from d using order.
func NewBinaryReader(d []byte, order binary.ByteOrder) *BinaryReader {
	return &BinaryReader{
		buf:   bytes.NewBuffer(d),
		order: order,
	}
}

// ReadString reads a null terminated string from the internal buffer and returns it.
func (r *BinaryReader) ReadString() (s string, err error) {
	b, err := r.buf.ReadBytes(0)
	if err != nil {
		return "", err
	} else if len(b) == 0 {
		return "", nil
	}

	return string(b[0 : len(b)-1]), nil
}

// Read reads structure data from the buffer into data.
// Data must be a pointer to a fixed-size value or a slice of fixed-size values.
// For more information see binary.Read.
func (r *BinaryReader) Read(data interface{}) error {
	return binary.Read(r.buf, r.order, data)
}
