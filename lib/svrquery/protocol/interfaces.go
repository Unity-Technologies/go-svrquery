package protocol

import (
	"io"

	"github.com/netdata/go-orchestrator/module"
)

// Queryer is an interface implemented by all svrquery protocols.
type Queryer interface {
	Query() (Responser, error)
}

// Responser is an interface implemented by types which represent a query response.
type Responser interface {
	NumClients() int64
	MaxClients() int64
}

// Mapper represents something which can return the current map.
type Mapper interface {
	Map() string
}

// Client is an interface which is implemented by types which can act a query transport.
type Client interface {
	io.ReadWriteCloser
	Key() string
	Address() string
}

// Charter is an interface which is implemented by types which support custom netdata
// charts.
type Charter interface {
	Charts(serverID int64) module.Charts
}

// Collector is an interface which is implemented by Responsers that provide custom
// netdata metrics.
type Collector interface {
	Collect(serverID int64, mx map[string]int64)
}
