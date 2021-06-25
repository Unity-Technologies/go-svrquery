package sqp

import "github.com/multiplay/go-svrquery/lib/svrsample/protocol"

func init() {
	protocol.MustRegister("sqp", NewQueryResponder)
}
