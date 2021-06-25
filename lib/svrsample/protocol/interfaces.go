package protocol

// QueryResponder represents an interface to a concrete type which responds
// to query requests.
type QueryResponder interface {
	Respond(clientAddress string, buf []byte) ([]byte, error)
}
