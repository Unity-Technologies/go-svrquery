package sqp

const (
	// TODO(steve): remove this?
	// DefaultMaxPacketSize is the default maximum size of a packet (MTU 1500 - UDP+IP header size)
	DefaultMaxPacketSize = 1472

	// Version is the query protocol version this client uses.
	Version = uint16(1)
)
