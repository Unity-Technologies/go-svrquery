package sqp

import (
	"github.com/multiplay/go-svrquery/lib/svrquery/protocol"
)

func init() {
	protocol.MustRegister(protocol.SQP, newCreator)
}
