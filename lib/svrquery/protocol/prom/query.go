package prom

import (
	"fmt"
	"github.com/multiplay/go-svrquery/lib/svrquery"
	"strconv"

	"github.com/multiplay/go-svrquery/lib/svrquery/protocol"
	"github.com/prometheus/common/expfmt"
)

type queryer struct {
	c protocol.Client
}

func newCreator(c protocol.Client) protocol.Queryer {
	return newQueryer(c)
}

func newQueryer(c protocol.Client) *queryer {
	return &queryer{
		c: c,
	}
}

// Query implements protocol.Queryer.
func (q *queryer) Query() (protocol.Responser, error) {
	return q.makeQuery()
}

func (q *queryer) makeQuery() (*QueryResponse, error) {
	client, ok := q.c.(*svrquery.Client)
	if !ok {
		return nil, fmt.Errorf("expected svrquery.Client, got %T", q.c)
	}
	httpTransport, ok := client.Transport.(svrquery.HTTPTransport)
	if !ok {
		return nil, fmt.Errorf("expected svrquery.HTTPTransport, got %T", q.c)
	}
	res, err := httpTransport.HttpClient.Get(q.c.Address())
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	var parser expfmt.TextParser
	metrics, err := parser.TextToMetricFamilies(res.Body)
	if err != nil {
		return nil, err
	}

	resp := &QueryResponse{}
	for _, v := range metrics {
		switch *v.Name {
		case currentPlayersMetricName:
			resp.CurrentPlayers = *v.Metric[0].Gauge.Value
		case maxPlayersMetricName:
			resp.MaxPlayers = *v.Metric[0].Gauge.Value
		case serverInfoMetricName:
			if len(v.Metric) == 0 || v.Metric[0] == nil || len(v.Metric[0].Label) == 0 {
				// server_info metric is missing labels
				continue
			}
			for _, l := range v.Metric[0].Label {
				switch *l.Name {
				case "server_name":
					resp.ServerName = *l.Value
				case "game_type":
					resp.GameType = *l.Value
				case "map_name":
					resp.MapName = *l.Value
				case "port":
					portInt, err := strconv.ParseInt(*l.Value, 10, 64)
					if err != nil {
						// invalid port
						break
					}
					resp.Port = portInt
				}
			}
		}
	}

	return resp, nil
}
