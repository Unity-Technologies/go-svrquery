package prom

import (
	"github.com/multiplay/go-svrquery/lib/svrsample/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// QueryResponder responds to queries
type QueryResponder struct {
	state       common.QueryState
	registry    *prometheus.Registry
	metrics     *metrics
	HTTPHandler http.Handler
}

// NewQueryResponder returns creates a new responder capable of returning metrics in Prometheus format.
func NewQueryResponder(state common.QueryState) (*QueryResponder, error) {
	// Create a registry, new metrics and register them
	reg := prometheus.NewRegistry()
	m := newMetrics(reg)

	httpHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg})

	q := &QueryResponder{
		state:       state,
		registry:    reg,
		metrics:     m,
		HTTPHandler: httpHandler,
	}

	// update metrics (though in this sample server they never change)
	m.UpdateFromQueryState(state)

	return q, nil
}

// Respond generates the query response for the requester.
func (q *QueryResponder) Respond(_ string, _ []byte) ([]byte, error) {
	// no-op, the http handler takes care of writing responses
	return nil, nil
}
