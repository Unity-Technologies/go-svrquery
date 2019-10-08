package titanfall

const (
	// ServerInfoRequest is the command type of a info request packet
	ServerInfoRequest = byte(77)

	// ServerInfoResponse is the command type of a info response packet
	ServerInfoResponse = byte(78)

	// ServerInfoVersion is the version of a info packets.
	ServerInfoVersion = byte(3)

	// ServerInfoVersionKeyed is the version of keys info packets.
	ServerInfoVersionKeyed = byte(5)
)
