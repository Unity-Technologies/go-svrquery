package titanfall

import (
	"github.com/multiplay/go-svrquery/lib/svrquery/protocol"
)

func init() {
	// TODO(steve): add support for tf2.
	protocol.MustRegister(protocol.TF2, newQueryer(3))
	protocol.MustRegister(protocol.TF2v7, newQueryer(7))
	protocol.MustRegister(protocol.TF2v8, newQueryer(8))
	protocol.MustRegister(protocol.TF2v9, newQueryer(9))
	protocol.MustRegister(protocol.TF2v10, newQueryer(10))
}
