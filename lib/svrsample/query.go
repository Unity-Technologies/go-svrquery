package svrsample

import (
	"errors"
	"fmt"

	"github.com/multiplay/go-svrquery/lib/svrsample/common"

	"github.com/multiplay/go-svrquery/lib/svrsample/protocol/sqp"
)

var (
	// ErrProtoNotFound returned when a protocol is not found
	ErrProtoNotFound = errors.New("protocol not found")
)

// GetResponder gets the appropriate responder for the protocol provided
func GetResponder(proto string, state common.QueryState) (common.QueryResponder, error) {
	switch proto {
	case "sqp":
		return sqp.NewQueryResponder(state)
	}
	return nil, fmt.Errorf("%w: %s", ErrProtoNotFound, proto)
}
