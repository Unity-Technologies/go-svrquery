package titanfall

import (
	"github.com/multiplay/go-svrquery/lib/svrquery/protocol"
)

func init() {
	// TODO(steve): add support for tf2.
	protocol.MustRegister("tf2e", newQueryer)
}
