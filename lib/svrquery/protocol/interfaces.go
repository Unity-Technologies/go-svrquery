package protocol

import (
	"io"
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
	Queryer
	Key() string
	Address() string
}
