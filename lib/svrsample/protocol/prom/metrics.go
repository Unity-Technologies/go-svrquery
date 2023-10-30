package prom

import (
	"github.com/multiplay/go-svrquery/lib/svrsample/common"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

var _ prometheus.Collector = (*metrics)(nil)

const (
	metricNamespace = "gameserver"
)

// metrics holds the current prometheus metrics data for the server
type metrics struct {
	currentPlayers prometheus.Gauge
	maxPlayers     prometheus.Gauge
	serverInfo     *prometheus.GaugeVec
}

func newMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		currentPlayers: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: metricNamespace,
			Subsystem: "",
			Name:      "current_players",
			Help:      "Number of players currently connected to the server.",
		}),
		maxPlayers: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: metricNamespace,
			Subsystem: "",
			Name:      "max_players",
			Help:      "Maximum number of players that can connect to the server.",
		}),
		serverInfo: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace:   metricNamespace,
			Subsystem:   "",
			Name:        "server_info",
			Help:        "Server status info.",
			ConstLabels: nil,
		}, []string{"server_name", "game_type", "map_name", "port"}),
	}
	reg.MustRegister(m)
	return m
}

func (m metrics) Describe(descs chan<- *prometheus.Desc) {
	m.currentPlayers.Describe(descs)
	m.maxPlayers.Describe(descs)
	m.serverInfo.Describe(descs)
}

func (m metrics) Collect(c chan<- prometheus.Metric) {
	m.currentPlayers.Collect(c)
	m.maxPlayers.Collect(c)
	m.serverInfo.Collect(c)
}

// UpdateFromQueryState populates metrics with data from common.QueryState.
func (m metrics) UpdateFromQueryState(qs common.QueryState) {
	m.currentPlayers.Set(float64(qs.CurrentPlayers))
	m.maxPlayers.Set(float64(qs.MaxPlayers))
	portString := strconv.FormatUint(uint64(qs.Port), 10)
	m.serverInfo.WithLabelValues(qs.ServerName, qs.GameType, qs.Map, portString).Set(1)
}
