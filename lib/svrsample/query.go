package svrsample

import (
	"errors"
	"fmt"

	"github.com/multiplay/go-svrquery/lib/svrsample/protocol/sqp"
)

var (
	// ErrProtoNotFound returned when a protocol is not found
	ErrProtoNotFound = errors.New("protocol not found")
)

// QueryResponder represents an interface to a concrete type which responds
// to query requests.
type QueryResponder interface {
	Respond(clientAddress string, buf []byte) ([]byte, error)
}

// QueryState represents the state of a currently running game.
type QueryState struct {
	CurrentPlayers int32
	MaxPlayers     int32
	ServerName     string
	GameType       string
	Map            string
	Port           uint16
}

// GetResponder gets the appropriate responder for the protocol provided
func GetResponder(proto string, state QueryState) (QueryResponder, error) {
	switch proto {
	case "sqp":
		return sqp.NewQueryResponder(state)
	}
	return nil, fmt.Errorf("%w: %s", ErrProtoNotFound, proto)
}
