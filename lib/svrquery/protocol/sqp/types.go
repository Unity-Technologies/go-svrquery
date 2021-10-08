package sqp

import (
	"encoding/json"
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

// QueryResponse is the combined response to a query request
type QueryResponse struct {
	Version     uint16            `json:"version"`
	Address     string            `json:"address"`
	ServerInfo  *ServerInfoChunk  `json:"server_info,omitempty"`
	ServerRules *ServerRulesChunk `json:"server_rules,omitempty"`
	PlayerInfo  *PlayerInfoChunk  `json:"player_info,omitempty"`
	TeamInfo    *TeamInfoChunk    `json:"team_info,omitempty"`
}

// MaxClients returns the maximum number of clients.
func (q *QueryResponse) MaxClients() int64 {
	if q.ServerInfo == nil {
		// No server info chunk, use 0
		return 0
	}
	return int64(q.ServerInfo.MaxPlayers)
}

// NumClients returns the number of clients.
func (q *QueryResponse) NumClients() int64 {
	if q.ServerInfo == nil {
		// No server info chunk, use 0
		return 0
	}
	return int64(q.ServerInfo.CurrentPlayers)
}

type infoHeader struct {
	Name string
	Type DataType
}
