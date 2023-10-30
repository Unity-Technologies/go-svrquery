package prom

const (
	currentPlayersMetricName = "current_players"
	maxPlayersMetricName     = "max_players"
	serverInfoMetricName     = "server_info"
)

// QueryResponse is the combined response to a query request
type QueryResponse struct {
	CurrentPlayers float64 `json:"current_players"`
	MaxPlayers     float64 `json:"max_players"`
	ServerName     string  `json:"server_name"`
	GameType       string  `json:"game_type"`
	MapName        string  `json:"map"`
	Port           int64   `json:"port"`
	//BuildID        string `json:"build_id"`
}

// MaxClients implements protocol.Responser, returns the maximum number of clients.
func (q *QueryResponse) MaxClients() int64 {
	return int64(q.MaxPlayers)
}

// NumClients implements protocol.Responser, returns the number of clients.
func (q *QueryResponse) NumClients() int64 {
	return int64(q.CurrentPlayers)
}

func (q *QueryResponse) Map() string {
	return q.MapName
}
