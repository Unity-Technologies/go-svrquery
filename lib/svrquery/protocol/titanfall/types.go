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
	// Version 9+
	PerformanceInfoV9
	// Version 2+
	MatchState
	// Version 9+
	MatchStateV9

	Teams   []Team
	Clients []Client
}

// NumClients implements protocol.Responser.
func (i Info) NumClients() int64 {
	return int64(i.BasicInfo.NumClients)
}

// NumBotClients implements protocol.Responser.
func (i Info) NumBotClients() int64 {
	return int64(i.BasicInfo.NumBotClients)
}

// MaxClients implements protocol.Responser.
func (i Info) MaxClients() int64 {
	return int64(i.BasicInfo.MaxClients)
}

// TotalClientsConnectedEver implements protocol.Responser.
func (i Info) TotalClientsConnectedEver() int64 {
	return int64(i.BasicInfo.TotalClientsConnectedEver)
}

// Map implements protocol.Mapper.
func (i Info) Map() string {
	return i.BasicInfo.Map
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

// InstanceInfoV8 represents instance information contained in a query response.
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
		DOS              bool
	}{}
	obj.None = a.None()
	obj.PacketLossIn = a.PacketLossIn()
	obj.PacketLossOut = a.PacketLossOut()
	obj.PacketChokedIn = a.PacketChokedIn()
	obj.PacketChokedOut = a.PacketChokedOut()
	obj.SlowServerFrames = a.SlowServerFrames()
	obj.Hitching = a.Hitching()
	obj.DOS = a.DOS()

	return json.Marshal(obj)
}

// None No health issues
func (a HealthFlags) None() bool {
	return a == 0
}

// PacketLossIn health flag
func (a HealthFlags) PacketLossIn() bool {
	return (a>>0)&1 == 1
}

// PacketLossOut health flag
func (a HealthFlags) PacketLossOut() bool {
	return (a>>1)&1 == 1
}

// PacketChokedIn health flag
func (a HealthFlags) PacketChokedIn() bool {
	return (a>>2)&1 == 1
}

// PacketChokedOut health flag
func (a HealthFlags) PacketChokedOut() bool {
	return (a>>3)&1 == 1
}

// SlowServerFrames health flag
func (a HealthFlags) SlowServerFrames() bool {
	return (a>>4)&1 == 1
}

// Hitching health flag
func (a HealthFlags) Hitching() bool {
	return (a>>5)&1 == 1
}

// DOS health flag
func (a HealthFlags) DOS() bool {
	return (a>>6)&1 == 1
}

// BasicInfo represents basic information contained in a query response.
type BasicInfo struct {
	Port                      uint16
	Platform                  string
	PlaylistVersion           string
	PlaylistNum               uint32
	PlaylistName              string
	NumClients                byte
	NumBotClients             byte
	MaxClients                byte
	TotalClientsConnectedEver uint32
	Map                       string
	PlatformPlayers           map[string]byte
}

// PerformanceInfoV9 represents frame information contained in a query response.
type PerformanceInfoV9 struct {
	PerformanceInfo
	CommitMemory   uint32
	ResidentMemory uint32
}

// PerformanceInfo represents frame information contained in a query response.
type PerformanceInfo struct {
	AverageFrameTime       float32
	MaxFrameTime           float32
	AverageUserCommandTime float32
	MaxUserCommandTime     float32
}

// MatchStateV9 represents match state contained in a query response.
// A number of fields were removed in this version of the schema.
type MatchStateV9 struct {
	Phase      byte
	TimePassed uint16 // seconds
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
