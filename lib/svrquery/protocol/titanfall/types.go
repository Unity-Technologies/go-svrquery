package titanfall

import (
	"encoding/json"
	"fmt"

	"github.com/multiplay/go-svrquery/lib/svrquery/common"
)

// Info represents a full query response.
type Info struct {
	// All
	Header
	// Version 1+
	InstanceInfo
	// Version 8+
	InstanceInfoV8
	BuildName  string
	Datacenter string
	GameMode   string
	// All
	BasicInfo
	// Version 4+
	PerformanceInfo
	// Version 2+
	MatchState
	Teams   []Team
	Clients []Client
}

// NumClients implements protocol.Responser.
func (i Info) NumClients() int64 {
	return int64(i.BasicInfo.NumClients)
}

// MaxClients implements protocol.Responser.
func (i Info) MaxClients() int64 {
	return int64(i.BasicInfo.MaxClients)
}

// Header represents the header of a query response.
type Header struct {
	Prefix  int32
	Command byte
	Version byte
}

// InstanceInfo represents instance information contained in a query response.
type InstanceInfo struct {
	Retail         byte
	InstanceType   byte
	ClientCRC      uint32
	NetProtocol    uint16
	RandomServerID uint64
}

// InstanceInfo represents instance information contained in a query response.
type InstanceInfoV8 struct {
	Retail         byte
	InstanceType   byte
	ClientCRC      uint32
	NetProtocol    uint16
	HealthFlags    HealthFlags
	RandomServerID uint32
}

// HealthFlags allows us to read the health bits and output them
// in an easy to consume json format.
type HealthFlags uint32

// MarshalJSON implements json.Marshaler
func (a HealthFlags) MarshalJSON() ([]byte, error) {
	obj := struct {
		None             bool
		PacketLossIn     bool
		PacketLossOut    bool
		PacketChokedIn   bool
		PacketChokedOut  bool
		SlowServerFrames bool
		Hitching         bool
	}{}
	obj.None = a == 0
	obj.PacketLossIn = (a>>0)&1 == 1
	obj.PacketLossOut = (a>>1)&1 == 1
	obj.PacketChokedIn = (a>>2)&1 == 1
	obj.PacketChokedOut = (a>>3)&1 == 1
	obj.SlowServerFrames = (a>>4)&1 == 1
	obj.Hitching = (a>>5)&1 == 1

	return json.Marshal(obj)
}

// BasicInfo represents basic information contained in a query response.
type BasicInfo struct {
	Port            uint16
	Platform        string
	PlaylistVersion string
	PlaylistNum     uint32
	PlaylistName    string
	NumClients      byte
	MaxClients      byte
	Map             string
	PlatformPlayers map[string]byte
}

// PerformanceInfo represents frame information contained in a query response.
type PerformanceInfo struct {
	AverageFrameTime       float32
	MaxFrameTime           float32
	AverageUserCommandTime float32
	MaxUserCommandTime     float32
}

// MatchStateV2 represents match state contained in a query response.
// This contains a legacy v2 version of matchstate
type MatchStateV2 struct {
	Phase            byte
	MaxRounds        byte
	RoundsWonIMC     byte
	RoundsWonMilitia byte
	TimeLimit        uint16 // seconds
	TimePassed       uint16 // seconds
	MaxScore         uint16
}

// MatchState represents match state contained in a query response.
type MatchState struct {
	MatchStateV2
	TeamsLeftWithPlayersNum uint16
}

// Team represents a team in a query response.
type Team struct {
	ID    byte
	Score uint16
}

// Client represents a team in a query response.
type Client struct {
	ID     uint64
	Name   string
	TeamID byte
	// Version 3+
	Address         string
	Ping            uint32
	PacketsReceived uint32
	PacketsDropped  uint32
	// Version 2+
	Score  uint32
	Kills  uint16
	Deaths uint16
}

// Collect implements protocol.Collector.
func (i Info) Collect(serverID int64, mx map[string]int64) {
	if i.Version >= 2 {
		mx[fmt.Sprintf("%d_phase", serverID)] = int64(i.Phase)
	}
	if i.Version >= 5 {
		mx[fmt.Sprintf("%d_avg_frame_time", serverID)] = int64(i.AverageFrameTime * common.Dim3DP)
		mx[fmt.Sprintf("%d_max_frame_time", serverID)] = int64(i.MaxFrameTime * common.Dim3DP)
		mx[fmt.Sprintf("%d_avg_user_cmd_time", serverID)] = int64(i.AverageUserCommandTime * common.Dim3DP)
		mx[fmt.Sprintf("%d_max_user_cmd_time", serverID)] = int64(i.MaxUserCommandTime * common.Dim3DP)
	}
}
