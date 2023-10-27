package prom

import (
	"bytes"
	"github.com/multiplay/go-svrquery/lib/svrsample/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// QueryResponder responds to queries
type QueryResponder struct {
	state   common.QueryState
	metrics *metrics
}

type metrics struct {
	cpuTemp    prometheus.Gauge
	hdFailures *prometheus.CounterVec
}

func newMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		cpuTemp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_temperature_celsius",
			Help: "Current temperature of the CPU.",
		}),
		hdFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "hd_errors_total",
				Help: "Number of hard-disk errors.",
			},
			[]string{"device"},
		),
	}
	reg.MustRegister(m.cpuTemp)
	reg.MustRegister(m.hdFailures)
	return m
}

// NewQueryResponder returns creates a new responder capable of responding
// to SQP-formatted queries.
func NewQueryResponder(state common.QueryState) (*QueryResponder, error) {
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	m := newMetrics(reg)

	// Expose metrics and custom registry via an HTTP server
	// using the HandleFor function. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":8080", nil))
	q := &QueryResponder{
		state:   state,
		metrics: m,
	}
	return q, nil
}

// Respond writes a query response to the requester in the SQP wire protocol.
func (q *QueryResponder) Respond(_ string, buf []byte) ([]byte, error) {
	return q.handleQuery(buf)
}

// handleQuery handles an incoming query packet.
func (q *QueryResponder) handleQuery(buf []byte) ([]byte, error) {
	// Set values for the new created metrics.
	m.cpuTemp.Set(65.3)
	m.hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()

	wantsServerInfo := requestedChunks&0x1 == 1
	wantsMetrics := requestedChunks&0x10 == 16

	f := queryWireFormat{
		Header:     1,
		Challenge:  expectedChallenge.(uint32),
		SQPVersion: 1,
	}

	resp := bytes.NewBuffer(nil)

	if wantsServerInfo {
		f.ServerInfo = ServerInfoFromQueryState(q.state)
		size := f.ServerInfo.Size()
		f.ServerInfoLength = &size
		f.PayloadLength += uint16(*f.ServerInfoLength) + 4
	}

	if wantsMetrics {
		f.Metrics = MetricsFromQueryState(q.state)
		size := f.Metrics.Size()
		f.MetricsLength = &size
		f.PayloadLength += uint16(*f.MetricsLength) + 4
	}

	if err := common.WireWrite(resp, q.enc, f); err != nil {
		return nil, err
	}

	return resp.Bytes(), nil
}
