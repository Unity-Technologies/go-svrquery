package sqp

import (
	"github.com/multiplay/go-svrquery/lib/svrquery/protocol"
)

func init() {
	protocol.MustRegister("sqp", newCreator(ServerInfo))
	protocol.MustRegister("sqp2", newCreator(ServerInfo|Metrics))
}
