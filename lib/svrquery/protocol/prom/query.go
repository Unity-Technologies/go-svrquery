package prom

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/multiplay/go-svrquery/lib/svrquery/protocol"
	"github.com/prometheus/common/expfmt"
)

const defaultBufSize = 4096

type queryer struct {
	client protocol.Client
}

func newCreator(c protocol.Client) protocol.Queryer {
	return newQueryer(c)
}

func newQueryer(client protocol.Client) *queryer {
	return &queryer{
		client: client,
	}
}

// Query implements protocol.Queryer.
func (q *queryer) Query() (protocol.Responser, error) {
	return q.makeQuery()
}

func (q *queryer) makeQuery() (*QueryResponse, error) {
	// FIXME: this won't work if the response is larger than defaultBufSize
	responseBytes := make([]byte, defaultBufSize)
	n, err := q.client.Read(responseBytes)
	if n > 0 {
		responseBytes = responseBytes[:n]
	}
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("query response: %w", err)
	}
	var parser expfmt.TextParser
	metrics, err := parser.TextToMetricFamilies(bytes.NewReader(responseBytes))
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
