package prom

import (
	"github.com/multiplay/go-svrquery/lib/svrsample/common"
)

// Metrics holds the Metrics chunk data.
type Metrics struct {
	Count  byte
	Values []float32
}

// Size returns the number of bytes QueryResponder will use on the wire.
func (m Metrics) Size() uint32 {
	return uint32(1 + len(m.Values)*4)
}

// MetricsFromQueryState converts metrics data in common.QueryState to Metrics.
func MetricsFromQueryState(qs common.QueryState) *Metrics {
	l := len(qs.Metrics)
	m := &Metrics{
		Count: byte(l),
	}

	m.Values = make([]float32, l)
	copy(m.Values, qs.Metrics)

	return m
}
