package sqp

//go:generate stringer -type=DataType -output=enums_string.go

// DataType indicates the type of a response field
type DataType byte

// Size returns the DataTypes size in bytes, or -1 if unknown
func (dt DataType) Size() int {
	if dt > Uint64 {
		return -1
	}
	return 1 << dt
}

// Supported types for response fields
const (
	Byte DataType = iota
	Uint16
	Uint32
	Uint64
	String
)

// Request Types
const (
	ChallengeRequestType byte = iota
	QueryRequestType
)

// Response Types
const (
	ChallengeResponseType byte = iota
	QueryResponseType
)

// Query Requested Chunks
const (
	ServerInfo byte = 1 << iota
	ServerRules
	PlayerInfo
	TeamInfo
	PerformanceInfo
)
