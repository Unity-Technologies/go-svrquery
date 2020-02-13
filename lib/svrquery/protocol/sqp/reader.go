package sqp

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
	"unicode/utf8"
)

// packetReader is a collection of helpers for reading
// parts of a packet
type packetReader struct {
	io.Reader
}

// newPacketReader returns a new packetReader
func newPacketReader(r io.Reader) *packetReader {
	return &packetReader{r}
}

// ReadUint16 returns a uint16 from the underlying reader
func (pr *packetReader) ReadUint16() (v uint16, err error) {
	return v, binary.Read(pr, binary.BigEndian, &v)
}

// ReadUint32 returns a uint32 from the underlying reader
func (pr *packetReader) ReadUint32() (v uint32, err error) {
	return v, binary.Read(pr, binary.BigEndian, &v)
}

// ReadUint64 returns a uint64 from the underlying reader
func (pr *packetReader) ReadUint64() (v uint64, err error) {
	return v, binary.Read(pr, binary.BigEndian, &v)
}

// ReadByte returns a byte from the underlying reader
func (pr *packetReader) ReadByte() (v byte, err error) {
	return v, binary.Read(pr, binary.BigEndian, &v)
}

// ReadString returns a string and the number of bytes representing it (len byte + len) from the underlying reader
func (pr *packetReader) ReadString() (int64, string, error) {
	// Read the first byte as the length of the string
	length, err := pr.ReadByte()
	if err != nil {
		return 0, "", err
	}

	// Get the actual string data
	buf := make([]byte, length)
	n, err := pr.Read(buf)
	if err != nil {
		return int64(n + 1), "", err
	} else if n != int(length) {
		return int64(n + 1), "", ErrMalformedPacket(fmt.Sprintf("readstring: expected string with length %v, but read %v bytes of string data", length, n))
	}

	// Check it is valid utf8
	if !utf8.Valid(buf) {
		return int64(n + 1), "", ErrInvalidString
	}

	return int64(length + 1), string(buf), err
}

// deadlineReadWriter is a reader that applies a connection deadline before every read
type deadlineReadWriter struct {
	timeout time.Duration
	net.Conn
}

// NewDeadlineReadWriter returns a new deadlineReadWriter
func newDeadlineReadWriter(c net.Conn, t time.Duration) *deadlineReadWriter {
	return &deadlineReadWriter{Conn: c, timeout: t}
}

// Read implements the io.Reader interface
func (dr *deadlineReadWriter) Read(p []byte) (int, error) {
	if err := dr.Conn.SetReadDeadline(time.Now().Add(dr.timeout)); err != nil {
		return 0, err
	}
	return dr.Conn.Read(p)
}

// Write implements the io.Writer interface
func (dr *deadlineReadWriter) Write(p []byte) (int, error) {
	if err := dr.Conn.SetWriteDeadline(time.Now().Add(dr.timeout)); err != nil {
		return 0, err
	}
	return dr.Conn.Write(p)
}
