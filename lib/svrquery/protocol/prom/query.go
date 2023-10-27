package prom

import (
	"fmt"
	"github.com/multiplay/go-svrquery/lib/svrquery/protocol"
	"github.com/prometheus/common/expfmt"
	"net/http"
)

type queryer struct {
	c          protocol.Client
	httpClient *http.Client
	//	maxPktSize      int
	//	reader          *packetReader
	//	challengeID     uint32
	//	requestedChunks byte
}

func newCreator(c protocol.Client) protocol.Queryer {
	return newQueryer(c)
}

func newQueryer(c protocol.Client) *queryer {
	return &queryer{
		c:          c,
		httpClient: &http.Client{},
	}
}

// Query implements protocol.Queryer.
func (q *queryer) Query() (protocol.Responser, error) {
	if err := q.sendQuery(); err != nil {
		return nil, err
	}

	//return q.readQuery(q.requestedChunks)
	return nil, nil
}

func (q *queryer) sendQuery() error {
	res, err := q.httpClient.Get(q.c.Address())
	if err != nil {
		return fmt.Errorf("http get: %w", err)
	}
	var parser expfmt.TextParser
	metrics, err := parser.TextToMetricFamilies(res.Body)
	if err != nil {
		return err
	}

	for _, v := range metrics {
		ms := v.Metric[0].Label
		// TODO: continue
	}

	return nil
}

//func (q *queryer) readQueryServerInfo(qr *QueryResponse, r *packetReader) (err error) {
//	qr.ServerInfo = &ServerInfoChunk{}
//
//	if qr.ServerInfo.ChunkLength, err = r.ReadUint32(); err != nil {
//		return err
//	}
//
//	l := int64(qr.ServerInfo.ChunkLength)
//	if qr.ServerInfo.CurrentPlayers, err = r.ReadUint16(); err != nil {
//		return err
//	}
//	l -= int64(Uint16.Size())
//
//	if qr.ServerInfo.MaxPlayers, err = r.ReadUint16(); err != nil {
//		return err
//	}
//	l -= int64(Uint16.Size())
//
//	var n int64
//	if n, qr.ServerInfo.ServerName, err = r.ReadString(); err != nil {
//		return err
//	}
//	l -= n
//
//	if n, qr.ServerInfo.GameType, err = r.ReadString(); err != nil {
//		return err
//	}
//	l -= n
//
//	if n, qr.ServerInfo.BuildID, err = r.ReadString(); err != nil {
//		return err
//	}
//	l -= n
//
//	if n, qr.ServerInfo.Map, err = r.ReadString(); err != nil {
//		return err
//	}
//	l -= n
//
//	if qr.ServerInfo.Port, err = r.ReadUint16(); err != nil {
//		return err
//	}
//	l -= int64(Uint16.Size())
//
//	if l < 0 {
//		// If we have read more bytes than expected, the packet is malformed
//		return NewErrMalformedPacketf("expected chunk length of %v, but have %v bytes remaining", qr.ServerInfo.ChunkLength, l)
//	} else if l > 0 {
//		// If we have extra bytes remaining, we assume they are new fields from a future
//		// query version and discard them
//		if _, err := io.CopyN(ioutil.Discard, r, l); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func (q *queryer) readQueryServerRules(qr *QueryResponse, r *packetReader) (err error) {
//	qr.ServerRules = &ServerRulesChunk{Rules: make(map[string]*DynamicValue)}
//
//	if qr.ServerRules.ChunkLength, err = r.ReadUint32(); err != nil {
//		return err
//	}
//
//	l := int64(qr.ServerRules.ChunkLength)
//
//	for l > 0 {
//		n, name, err := r.ReadString()
//		if err != nil {
//			return err
//		}
//		l -= n
//
//		n, qr.ServerRules.Rules[name], err = NewDynamicValue(r)
//		if err != nil {
//			return err
//		}
//		l -= n
//	}
//
//	if l < 0 {
//		// If we have read more bytes than expected, the packet is malformed
//		return NewErrMalformedPacketf("expected chunk length of %v, but have %v bytes remaining", qr.ServerRules.ChunkLength, l)
//	}
//
//	return nil
//}
//
//func (q *queryer) readInfoHeader(r *packetReader) (int64, []*infoHeader, error) {
//	var n int64
//
//	expectedFieldCount, err := r.ReadByte()
//	if err != nil {
//		return 0, nil, err
//	}
//	n++
//
//	// If there are no fields, error as we should not have had records if there were no fields
//	if expectedFieldCount == 0 {
//		return n, nil, ErrMalformedPacket("no fields in info header")
//	}
//
//	// Build up header info
//	header := make([]*infoHeader, 0, expectedFieldCount)
//
//	for expectedFieldCount > 0 {
//		c, name, err := r.ReadString()
//		if err != nil {
//			return 0, nil, err
//		}
//		n += c
//
//		dt, err := r.ReadByte()
//		if err != nil {
//			return 0, nil, err
//		}
//		n++
//
//		header = append(header, &infoHeader{
//			Name: name,
//			Type: DataType(dt),
//		})
//
//		expectedFieldCount--
//	}
//
//	return n, header, nil
//}
//
//func (q *queryer) readQueryPlayerInfo(qr *QueryResponse, r *packetReader) (err error) {
//	qr.PlayerInfo = &PlayerInfoChunk{}
//
//	if qr.PlayerInfo.ChunkLength, err = r.ReadUint32(); err != nil {
//		return err
//	}
//
//	l := int64(qr.PlayerInfo.ChunkLength)
//	expectedPlayerCount, err := r.ReadUint16()
//	if err != nil {
//		return err
//	}
//	l -= int64(Uint16.Size())
//
//	// If there are no players, just skip the whole chunk
//	if expectedPlayerCount == 0 {
//		if _, err = io.CopyN(ioutil.Discard, r, l); err != nil {
//			return err
//		}
//
//		return nil
//	}
//
//	// Read the player fields header
//	n, header, err := q.readInfoHeader(r)
//	if err != nil {
//		return nil
//	}
//	l -= n
//
//	// Build the map of values for each player from the header
//	qr.PlayerInfo.Players = make([]map[string]*DynamicValue, expectedPlayerCount)
//	for i := 0; expectedPlayerCount > 0 && l > 0; i++ {
//		qr.PlayerInfo.Players[i] = make(map[string]*DynamicValue)
//		for _, ih := range header {
//			n, qr.PlayerInfo.Players[i][ih.Name], err = NewDynamicValueWithType(r, ih.Type)
//			if err != nil {
//				return err
//			}
//			l -= n
//		}
//		expectedPlayerCount--
//	}
//
//	switch {
//	case l < 0:
//		// If we have read more bytes than expected, the packet is malformed
//		return NewErrMalformedPacketf("expected chunk length of %v, but have %v bytes remaining", qr.PlayerInfo.ChunkLength, l)
//	case l > 0:
//		// If we have extra bytes remaining, we assume they are new fields from a future
//		// query version and discard them
//		if _, err := io.CopyN(ioutil.Discard, r, l); err != nil {
//			return err
//		}
//	case expectedPlayerCount != 0:
//		return NewErrMalformedPacketf("expected %v player records, but got %v", len(qr.PlayerInfo.Players)+int(expectedPlayerCount), len(qr.PlayerInfo.Players))
//	}
//
//	return nil
//}
//
//func (q *queryer) readQueryTeamInfo(qr *QueryResponse, r *packetReader) (err error) {
//	qr.TeamInfo = &TeamInfoChunk{}
//
//	if qr.TeamInfo.ChunkLength, err = r.ReadUint32(); err != nil {
//		return err
//	}
//
//	l := int64(qr.TeamInfo.ChunkLength)
//	expectedTeamCount, err := r.ReadUint16()
//	if err != nil {
//		return err
//	}
//	l -= int64(Uint16.Size())
//
//	// If there are no teams, just skip the whole chunk
//	if expectedTeamCount == 0 {
//		if _, err = io.CopyN(ioutil.Discard, r, l); err != nil {
//			return err
//		}
//
//		return nil
//	}
//
//	// Read the team fields header
//	n, header, err := q.readInfoHeader(r)
//	if err != nil {
//		return nil
//	}
//	l -= n
//
//	// Build the map of values for each team from the header
//	qr.TeamInfo.Teams = make([]map[string]*DynamicValue, expectedTeamCount)
//	for i := 0; expectedTeamCount > 0 && l > 0; i++ {
//		qr.TeamInfo.Teams[i] = make(map[string]*DynamicValue)
//		for _, ih := range header {
//			n, qr.TeamInfo.Teams[i][ih.Name], err = NewDynamicValueWithType(r, ih.Type)
//			if err != nil {
//				return err
//			}
//			l -= n
//		}
//		expectedTeamCount--
//	}
//
//	switch {
//	case l < 0:
//		// If we have read more bytes than expected, the packet is malformed
//		return NewErrMalformedPacketf("expected chunk length of %v, but have %v bytes remaining", qr.TeamInfo.ChunkLength, l)
//	case l > 0:
//		// If we have extra bytes remaining, we assume they are new fields from a future
//		// query version and discard them
//		if _, err := io.CopyN(ioutil.Discard, r, l); err != nil {
//			return err
//		}
//	case expectedTeamCount != 0:
//		return NewErrMalformedPacketf("expected %v Team records, but got %v", len(qr.TeamInfo.Teams)+int(expectedTeamCount), len(qr.TeamInfo.Teams))
//	}
//
//	return nil
//}
//
//func (q *queryer) readQueryMetrics(qr *QueryResponse, r *packetReader) (err error) {
//	qr.Metrics = &MetricsChunk{}
//
//	if qr.Metrics.ChunkLength, err = r.ReadUint32(); err != nil {
//		return err
//	}
//	l := int64(qr.Metrics.ChunkLength)
//
//	qr.Metrics.MetricCount, err = r.ReadByte()
//	if err != nil {
//		return err
//	}
//	l -= int64(Byte.Size())
//
//	if qr.Metrics.MetricCount > MaxMetrics {
//		return NewErrMalformedPacketf("metric count cannot be greater than %v, but got %v", MaxMetrics, qr.Metrics.MetricCount)
//	}
//
//	qr.Metrics.Metrics = make([]float32, qr.Metrics.MetricCount)
//	for i := 0; i < int(qr.Metrics.MetricCount) && l > 0; i++ {
//		qr.Metrics.Metrics[i], err = r.ReadFloat32()
//		if err != nil {
//			return err
//		}
//		l -= int64(Float32.Size())
//	}
//
//	if l < 0 {
//		// If we have read more bytes than expected, the packet is malformed
//		return NewErrMalformedPacketf("expected chunk length of %v, but have %v bytes remaining", qr.Metrics.ChunkLength, l)
//	} else if l > 0 {
//		// If we have extra bytes remaining, we assume they are new fields from a future
//		// query version and discard them
//		if _, err = io.CopyN(ioutil.Discard, r, l); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//*/
