package sqp

import (
	"encoding/json"
	"strconv"
)

// ServerInfoChunk is the response chunk for server info data
type ServerInfoChunk struct {
	ChunkLength    uint32 `json:"-"`
	CurrentPlayers uint16 `json:"current_players"`
	MaxPlayers     uint16 `json:"max_players"`
	ServerName     string `json:"server_name"`
	GameType       string `json:"game_type"`
	BuildID        string `json:"build_id"`
	Map            string `json:"map"`
	Port           uint16 `json:"port"`
}

// ServerRulesChunk is the response chunk for server rules data
type ServerRulesChunk struct {
	ChunkLength uint32 `json:"-"`
	Rules       map[string]*DynamicValue
}

// MarshalJSON returns the JSON representation of the server rules
func (src *ServerRulesChunk) MarshalJSON() ([]byte, error) {
	return json.Marshal(src.Rules)
}

// PlayerInfoChunk is the response chunk for player data
type PlayerInfoChunk struct {
	ChunkLength uint32 `json:"-"`
	Players     []map[string]*DynamicValue
}

// MarshalJSON returns the JSON representation of the player info
func (pic *PlayerInfoChunk) MarshalJSON() ([]byte, error) {
	return json.Marshal(pic.Players)
}

// TeamInfoChunk is the response chunk for team data
type TeamInfoChunk struct {
	ChunkLength uint32 `json:"-"`
	Teams       []map[string]*DynamicValue
}

// MarshalJSON returns the JSON representation of the team info
func (tic *TeamInfoChunk) MarshalJSON() ([]byte, error) {
	return json.Marshal(tic.Teams)
}

type PerformanceInfoChunk struct {
	ChunkLength uint32    `json:"-"`
	NumFlags    byte      `json:"-"`
	Flags       uint32    `json:"-"`
	Gauges      []float32 `json:"-"`
}

// MarshalJSON implements json.Marshaler
func (pi PerformanceInfoChunk) MarshalJSON() ([]byte, error) {
	obj := make(map[string]interface{}, pi.NumFlags)
	for i := 0; i < int(pi.NumFlags); i++ {
		obj["flag_"+strconv.Itoa(i)] = (pi.Flags>>i)&1 == 1
	}

	for i, f := range pi.Gauges {
		obj["gauge_"+strconv.Itoa(i)] = f
	}

	//	obj.PacketLossIn = a.PacketLossIn()
	//obj.PacketLossOut = a.PacketLossOut()
	//obj.PacketChokedIn = a.PacketChokedIn()
	//obj.PacketChokedOut = a.PacketChokedOut()
	//obj.SlowServerFrames = a.SlowServerFrames()
	//obj.Hitching = a.Hitching()

	return json.Marshal(obj)
}

// QueryResponse is the combined response to a query request
type QueryResponse struct {
	Version         uint16                `json:"version"`
	Address         string                `json:"address"`
	ServerInfo      *ServerInfoChunk      `json:"server_info,omitempty"`
	ServerRules     *ServerRulesChunk     `json:"server_rules,omitempty"`
	PlayerInfo      *PlayerInfoChunk      `json:"player_info,omitempty"`
	TeamInfo        *TeamInfoChunk        `json:"team_info,omitempty"`
	PerformanceInfo *PerformanceInfoChunk `json:"performance_info,omitempty"`
}

func (q *QueryResponse) MaxClients() int64 {
	return int64(q.ServerInfo.MaxPlayers)
}

func (q *QueryResponse) NumClients() int64 {
	return int64(q.ServerInfo.CurrentPlayers)
}

type infoHeader struct {
	Name string
	Type DataType
}
