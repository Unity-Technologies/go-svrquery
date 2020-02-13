package protocol

import (
	"fmt"
)

// Creator is a function which returns a Queryer.
type Creator func(c Client) Queryer

var (
	registry = make(map[string]Creator)
)

// MustRegister registers a protocol.
// Panics if the name is a duplicate.
func MustRegister(name string, f Creator) {
	if _, ok := registry[name]; ok {
		panic(fmt.Sprintf("%s is already in registry", name))
	}
	registry[name] = f
}

// Get returns the creator a protocol.
func Get(name string) (Creator, error) {
	f, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown protocol %q", name)
	}
	return f, nil
}

// Supported returns true if protocol name is supported.
func Supported(name string) bool {
	_, ok := registry[name]
	return ok
}
