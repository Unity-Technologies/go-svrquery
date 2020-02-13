package sqp

import (
	"errors"
	"fmt"
)

// ErrMalformedPacket is raised when a malformed packet is encountered
type ErrMalformedPacket string

func (e ErrMalformedPacket) Error() string {
	return fmt.Sprintf("malformed packet: %v", string(e))
}

// NewErrMalformedPacketf makes a new ErrMalformedPacket with the formatted string
func NewErrMalformedPacketf(format string, args ...interface{}) ErrMalformedPacket {
	return ErrMalformedPacket(fmt.Sprintf(format, args...))
}

var (
	// ErrMissingAddress is raised when no connection address is specified
	ErrMissingAddress = errors.New("missing address parameter")
	// ErrNilOption is raised when an option is a nil function pointer
	ErrNilOption = errors.New("options should not be nil")
	// ErrInvalidString is raised when a string is read and it is not valid utf8
	ErrInvalidString = errors.New("string is not valid utf8")
)

// ErrUnknownDataType is raised when a data type for a dynamic value is not recognised
type ErrUnknownDataType DataType

func (e ErrUnknownDataType) Error() string {
	return fmt.Sprintf("unknown datatype %v", string(e))
}
