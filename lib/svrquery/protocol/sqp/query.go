package sqp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"

	"github.com/multiplay/go-svrquery/lib/svrquery/protocol"
)

type queryer struct {
	c               protocol.Client
	maxPktSize      int
	reader          *packetReader
	challengeID     uint32
	requestedChunks byte
}

func newCreator(c protocol.Client) protocol.Queryer {
	return newQueryer(ServerInfo, DefaultMaxPacketSize, c)
}

func newQueryer(requestedChunks byte, maxPktSize int, c protocol.Client) *queryer {
	return &queryer{
		c:               c,
		maxPktSize:      maxPktSize,
		requestedChunks: requestedChunks,
		reader:          newPacketReader(bufio.NewReaderSize(c, maxPktSize)),
	}
}

// Query implements protocol.Queryer.
func (q *queryer) Query() (protocol.Responser, error) {
	if q.challengeID == 0 {
		if err := q.Challenge(); err != nil {
			return nil, err
		}
	}

	if err := q.sendQuery(q.requestedChunks); err != nil {
		return nil, err
	}

	return q.readQuery(q.requestedChunks)
}

func (q *queryer) sendQuery(requestedChunks byte) error {
	pkt := &bytes.Buffer{}
	if err := binary.Write(pkt, binary.BigEndian, QueryRequestType); err != nil {
		return err
	}

	if err := binary.Write(pkt, binary.BigEndian, q.challengeID); err != nil {
		return err
	}

	if err := binary.Write(pkt, binary.BigEndian, Version); err != nil {
		return err
	}

	if err := binary.Write(pkt, binary.BigEndian, requestedChunks); err != nil {
		return err
	}

	_, err := q.c.Write(pkt.Bytes())
	return err
}

func (q *queryer) readQueryHeader() (uint16, byte, byte, uint16, error) {
	pktType, err := q.reader.ReadByte()
	if err != nil {
		return 0, 0, 0, 0, err
	} else if pktType != QueryResponseType {
		return 0, 0, 0, 0, NewErrMalformedPacketf("was expecting %v for response type, got %v", QueryResponseType, pktType)
	}

	if err = q.readChallenge(); err != nil {
		return 0, 0, 0, 0, err
	}

	var version uint16
	if version, err = q.reader.ReadUint16(); err != nil {
		return 0, 0, 0, 0, err
	}

	var curPkt, lastPkt byte
	curPkt, err = q.reader.ReadByte()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	lastPkt, err = q.reader.ReadByte()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	pktLen, err := q.reader.ReadUint16()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	if curPkt > lastPkt {
		return 0, 0, 0, 0, ErrMalformedPacket("current packet id > last packet id")
	}

	return version, curPkt, lastPkt, pktLen, nil
}

func (q *queryer) readQuery(requestedChunks byte) (*QueryResponse, error) {
	version, curPkt, lastPkt, pktLen, err := q.readQueryHeader()
	if err != nil {
		return nil, err
	}

	if lastPkt == 0 && curPkt == 0 {
		// If the header says the body is empty, we should just return now
		if pktLen == 0 {
			return &QueryResponse{Version: version, Address: q.c.Address()}, nil
		}

		return q.readQuerySinglePacket(q.reader, version, requestedChunks, uint32(pktLen))
	}

	return q.readQueryMultiPacket(version, curPkt, lastPkt, requestedChunks, pktLen)
}

func (q *queryer) readQuerySinglePacket(r *packetReader, version uint16, requestedChunks byte, pktLen uint32) (*QueryResponse, error) {
	qr := &QueryResponse{Version: version, Address: q.c.Address()}

	l := pktLen
	if requestedChunks&ServerInfo > 0 {
		if err := q.readQueryServerInfo(qr, r); err != nil {
			return nil, err
		}
		l -= qr.ServerInfo.ChunkLength + uint32(Uint32.Size())
	}

	if requestedChunks&ServerRules > 0 {
		if err := q.readQueryServerRules(qr, r); err != nil {
			return nil, err
		}
		l -= qr.ServerRules.ChunkLength + uint32(Uint32.Size())
	}

	if requestedChunks&PlayerInfo > 0 {
		if err := q.readQueryPlayerInfo(qr, r); err != nil {
			return nil, err
		}
		l -= qr.PlayerInfo.ChunkLength + uint32(Uint32.Size())
	}

	if requestedChunks&TeamInfo > 0 {
		if err := q.readQueryTeamInfo(qr, r); err != nil {
			return nil, err
		}
		l -= qr.TeamInfo.ChunkLength + uint32(Uint32.Size())
	}

	if l < 0 {
		// If we have read more bytes than expected, the packet is malformed
		return nil, NewErrMalformedPacketf("expected packet length of %v, but have %v bytes remaining", pktLen, l)
	} else if l > 0 {
		// If we have extra bytes remaining, we assume they are new fields from a future
		// query version and discard them.
		if _, err := io.CopyN(ioutil.Discard, r, int64(l)); err != nil {
			return nil, err
		}
	}

	return qr, nil
}

func (q *queryer) readQueryServerInfo(qr *QueryResponse, r *packetReader) (err error) {
	qr.ServerInfo = &ServerInfoChunk{}

	if qr.ServerInfo.ChunkLength, err = r.ReadUint32(); err != nil {
		return err
	}

	l := int64(qr.ServerInfo.ChunkLength)
	if qr.ServerInfo.CurrentPlayers, err = r.ReadUint16(); err != nil {
		return err
	}
	l -= int64(Uint16.Size())

	if qr.ServerInfo.MaxPlayers, err = r.ReadUint16(); err != nil {
		return err
	}
	l -= int64(Uint16.Size())

	var n int64
	if n, qr.ServerInfo.ServerName, err = r.ReadString(); err != nil {
		return err
	}
	l -= n

	if n, qr.ServerInfo.GameType, err = r.ReadString(); err != nil {
		return err
	}
	l -= n

	if n, qr.ServerInfo.BuildID, err = r.ReadString(); err != nil {
		return err
	}
	l -= n

	if n, qr.ServerInfo.Map, err = r.ReadString(); err != nil {
		return err
	}
	l -= n

	if qr.ServerInfo.Port, err = r.ReadUint16(); err != nil {
		return err
	}
	l -= int64(Uint16.Size())

	if l < 0 {
		// If we have read more bytes than expected, the packet is malformed
		return NewErrMalformedPacketf("expected chunk length of %v, but have %v bytes remaining", qr.ServerInfo.ChunkLength, l)
	} else if l > 0 {
		// If we have extra bytes remaining, we assume they are new fields from a future
		// query version and discard them
		if _, err := io.CopyN(ioutil.Discard, r, l); err != nil {
			return err
		}
	}

	return nil
}

func (q *queryer) readQueryServerRules(qr *QueryResponse, r *packetReader) (err error) {
	qr.ServerRules = &ServerRulesChunk{Rules: make(map[string]*DynamicValue)}

	if qr.ServerRules.ChunkLength, err = r.ReadUint32(); err != nil {
		return err
	}

	l := int64(qr.ServerRules.ChunkLength)

	for l > 0 {
		n, name, err := r.ReadString()
		if err != nil {
			return err
		}
		l -= n

		n, qr.ServerRules.Rules[name], err = NewDynamicValue(r)
		if err != nil {
			return err
		}
		l -= n
	}

	if l < 0 {
		// If we have read more bytes than expected, the packet is malformed
		return NewErrMalformedPacketf("expected chunk length of %v, but have %v bytes remaining", qr.ServerRules.ChunkLength, l)
	}

	return nil
}

func (q *queryer) readInfoHeader(r *packetReader) (int64, []*infoHeader, error) {
	var n int64

	expectedFieldCount, err := r.ReadByte()
	if err != nil {
		return 0, nil, err
	}
	n++

	// If there are no fields, error as we should not have had records if there were no fields
	if expectedFieldCount == 0 {
		return n, nil, ErrMalformedPacket("no fields in info header")
	}

	// Build up header info
	header := make([]*infoHeader, 0, expectedFieldCount)

	for expectedFieldCount > 0 {
		c, name, err := r.ReadString()
		if err != nil {
			return 0, nil, err
		}
		n += c

		dt, err := r.ReadByte()
		if err != nil {
			return 0, nil, err
		}
		n++

		header = append(header, &infoHeader{
			Name: name,
			Type: DataType(dt),
		})

		expectedFieldCount--
	}

	return n, header, nil
}

func (q *queryer) readQueryPlayerInfo(qr *QueryResponse, r *packetReader) (err error) {
	qr.PlayerInfo = &PlayerInfoChunk{}

	if qr.PlayerInfo.ChunkLength, err = r.ReadUint32(); err != nil {
		return err
	}

	l := int64(qr.PlayerInfo.ChunkLength)
	expectedPlayerCount, err := r.ReadUint16()
	if err != nil {
		return err
	}
	l -= int64(Uint16.Size())

	// If there are no players, just skip the whole chunk
	if expectedPlayerCount == 0 {
		if _, err = io.CopyN(ioutil.Discard, r, l); err != nil {
			return err
		}

		return nil
	}

	// Read the player fields header
	n, header, err := q.readInfoHeader(r)
	if err != nil {
		return nil
	}
	l -= n

	// Build the map of values for each player from the header
	qr.PlayerInfo.Players = make([]map[string]*DynamicValue, expectedPlayerCount)
	for i := 0; expectedPlayerCount > 0 && l > 0; i++ {
		qr.PlayerInfo.Players[i] = make(map[string]*DynamicValue)
		for _, ih := range header {
			n, qr.PlayerInfo.Players[i][ih.Name], err = NewDynamicValueWithType(r, ih.Type)
			if err != nil {
				return err
			}
			l -= n
		}
		expectedPlayerCount--
	}

	switch {
	case l < 0:
		// If we have read more bytes than expected, the packet is malformed
		return NewErrMalformedPacketf("expected chunk length of %v, but have %v bytes remaining", qr.PlayerInfo.ChunkLength, l)
	case l > 0:
		// If we have extra bytes remaining, we assume they are new fields from a future
		// query version and discard them
		if _, err := io.CopyN(ioutil.Discard, r, l); err != nil {
			return err
		}
	case expectedPlayerCount != 0:
		return NewErrMalformedPacketf("expected %v player records, but got %v", len(qr.PlayerInfo.Players)+int(expectedPlayerCount), len(qr.PlayerInfo.Players))
	}

	return nil
}

func (q *queryer) readQueryTeamInfo(qr *QueryResponse, r *packetReader) (err error) {
	qr.TeamInfo = &TeamInfoChunk{}

	if qr.TeamInfo.ChunkLength, err = r.ReadUint32(); err != nil {
		return err
	}

	l := int64(qr.TeamInfo.ChunkLength)
	expectedTeamCount, err := r.ReadUint16()
	if err != nil {
		return err
	}
	l -= int64(Uint16.Size())

	// If there are no teams, just skip the whole chunk
	if expectedTeamCount == 0 {
		if _, err = io.CopyN(ioutil.Discard, r, l); err != nil {
			return err
		}

		return nil
	}

	// Read the team fields header
	n, header, err := q.readInfoHeader(r)
	if err != nil {
		return nil
	}
	l -= n

	// Build the map of values for each team from the header
	qr.TeamInfo.Teams = make([]map[string]*DynamicValue, expectedTeamCount)
	for i := 0; expectedTeamCount > 0 && l > 0; i++ {
		qr.TeamInfo.Teams[i] = make(map[string]*DynamicValue)
		for _, ih := range header {
			n, qr.TeamInfo.Teams[i][ih.Name], err = NewDynamicValueWithType(r, ih.Type)
			if err != nil {
				return err
			}
			l -= n
		}
		expectedTeamCount--
	}

	switch {
	case l < 0:
		// If we have read more bytes than expected, the packet is malformed
		return NewErrMalformedPacketf("expected chunk length of %v, but have %v bytes remaining", qr.TeamInfo.ChunkLength, l)
	case l > 0:
		// If we have extra bytes remaining, we assume they are new fields from a future
		// query version and discard them
		if _, err := io.CopyN(ioutil.Discard, r, l); err != nil {
			return err
		}
	case expectedTeamCount != 0:
		return NewErrMalformedPacketf("expected %v Team records, but got %v", len(qr.TeamInfo.Teams)+int(expectedTeamCount), len(qr.TeamInfo.Teams))
	}

	return nil
}

func (q *queryer) readQueryMultiPacket(version uint16, curPkt, lastPkt, requestedChunks byte, pktLen uint16) (*QueryResponse, error) {
	// Setup our array of packet bodies
	expectedPkts := lastPkt + 1
	multiPkt := make(map[byte][]byte, expectedPkts)
	totalPktLen := uint32(pktLen)

	// Handle this first packet
	multiPkt[curPkt] = make([]byte, pktLen)
	n, err := q.reader.Read(multiPkt[curPkt])
	if err != nil {
		return nil, err
	} else if uint16(n) != pktLen {
		return nil, NewErrMalformedPacketf("expected packet length of %v, but read %v bytes", pktLen, n)
	}

	// Remember the challengeID so that we can verify each packet we are reading is
	// part of this multi-packet response
	challengeID := q.challengeID

	// Handle each subsequent packet until we have all of the ones we need
	for len(multiPkt) != int(expectedPkts) {
		version, curPkt, lastPkt, pktLen, err = q.readQueryHeader()
		if err != nil {
			return nil, err
		}

		// If this packet isn't part of the multi-packet response we are expecting, discard it
		if q.challengeID != challengeID {
			if _, err := io.CopyN(ioutil.Discard, q.reader, int64(pktLen)); err != nil {
				return nil, err
			}
			continue
		}

		totalPktLen += uint32(pktLen)
		multiPkt[curPkt] = make([]byte, pktLen)
		n, err := q.reader.Read(multiPkt[curPkt])
		if err != nil {
			return nil, err
		} else if uint16(n) != pktLen {
			return nil, NewErrMalformedPacketf("expected packet length of %v, but read %v bytes", pktLen, n)
		}
	}

	// Now recombine the packets into the right order.
	// Sure, this could be more efficient by not using a map before
	// and instead just inserting the packets into the correct place
	// in a slice using splicing, but for now this is easier.
	buf := &bytes.Buffer{}
	for curPkt = 0; curPkt <= lastPkt; curPkt++ {
		buf.Write(multiPkt[curPkt])
	}

	return q.readQuerySinglePacket(newPacketReader(buf), version, requestedChunks, totalPktLen)
}
