package sqp

import (
	"testing"

	"github.com/multiplay/go-svrquery/lib/svrquery/clienttest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testDir = "testdata"
)

func newClient(requestedChunks byte) (*clienttest.MockClient, *queryer) {
	m := &clienttest.MockClient{}
	m.On("Address").Return("127.0.0.1:8000")
	c := newQueryer(requestedChunks, DefaultMaxPacketSize, m)

	// Act as if we've already issued a challenge/response
	c.challengeID = 1

	return m, c
}

func TestQuery(t *testing.T) {
	cases := []struct {
		name   string
		chunks byte
		multi  int
		f      func(t *testing.T, c *queryer)
	}{
		{
			name:   "info_single",
			chunks: ServerInfo,
			f:      testQueryInfoSinglePacket,
		},
		{
			name:   "info_single_malformed",
			chunks: ServerInfo,
			f:      testQueryServerInfoSinglePacketMalformed,
		},
		{
			name:   "info_multi",
			chunks: ServerInfo,
			multi:  2,
			f:      testQueryServerInfoMultiPacket,
		},
		{
			name:   "rules",
			chunks: ServerRules,
			f:      testQueryServerRulesSinglePacket,
		},
		{
			name:   "player",
			chunks: PlayerInfo,
			f:      testQueryPlayerInfoSinglePacket,
		},
		{
			name:   "team",
			chunks: TeamInfo,
			f:      testQueryTeamInfoSinglePacket,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := clienttest.LoadData(t, testDir, tc.name+"_request")

			m, c := newClient(tc.chunks)
			m.On("Write", req).Return(len(req), nil)

			if tc.multi > 0 {
				pkts := clienttest.LoadMultiData(t, tc.multi, testDir, tc.name+"_response")
				for _, resp := range pkts {
					m.On("Read", mock.AnythingOfType("[]uint8")).Return(resp, nil).Once()
				}
			} else {
				resp := clienttest.LoadData(t, testDir, tc.name+"_response")
				m.On("Read", mock.AnythingOfType("[]uint8")).Return(resp, nil)
			}

			tc.f(t, c)
		})
	}
}

func testQueryInfoSinglePacket(t *testing.T, c *queryer) {
	r, err := c.Query()
	require.NoError(t, err, "query request failed")

	qr := r.(*QueryResponse)
	require.Equal(t, uint32(256), c.challengeID, "expected correct challenge id")

	require.NotNil(t, qr, "expected query response")
	require.NotNil(t, qr.ServerInfo, "expected server info")

	require.Equal(t, uint16(5), qr.ServerInfo.CurrentPlayers)
	require.Equal(t, uint16(10), qr.ServerInfo.MaxPlayers)
	require.Equal(t, "my server", qr.ServerInfo.ServerName)
	require.Equal(t, "ctf", qr.ServerInfo.GameType)
	require.Equal(t, "1", qr.ServerInfo.BuildID)
	require.Equal(t, "map", qr.ServerInfo.Map)
	require.Equal(t, uint16(1025), qr.ServerInfo.Port)
}

func testQueryServerInfoSinglePacketMalformed(t *testing.T, c *queryer) {
	_, err := c.Query()
	require.Error(t, err, "query request should have failed")

	_, ok := err.(ErrMalformedPacket)
	require.Truef(t, ok, "expected malformed packet err, got: %v", err)
}

func testQueryServerInfoMultiPacket(t *testing.T, c *queryer) {
	r, err := c.Query()
	require.NoError(t, err, "query request should not have failed")
	qr := r.(*QueryResponse)

	require.Equal(t, uint32(256), c.challengeID, "expected correct challenge id")

	require.NotNil(t, qr, "expected query response")
	require.NotNil(t, qr.ServerInfo, "expected server info")

	require.Equal(t, uint16(5), qr.ServerInfo.CurrentPlayers)
	require.Equal(t, uint16(10), qr.ServerInfo.MaxPlayers)
	require.Equal(t, "my server", qr.ServerInfo.ServerName)
	require.Equal(t, "ctf", qr.ServerInfo.GameType)
	require.Equal(t, "1", qr.ServerInfo.BuildID)
	require.Equal(t, "map", qr.ServerInfo.Map)
	require.Equal(t, uint16(1025), qr.ServerInfo.Port)
}

func testQueryServerRulesSinglePacket(t *testing.T, c *queryer) {
	r, err := c.Query()
	require.NoError(t, err, "query request should not have failed")
	qr := r.(*QueryResponse)

	require.Equal(t, uint32(256), c.challengeID, "expected correct challenge id")

	require.NotNil(t, qr, "expected query response")
	require.NotNil(t, qr.ServerRules, "expected server rules")
	require.Len(t, qr.ServerRules.Rules, 5)

	require.Equal(t, byte(128), qr.ServerRules.Rules["rule 1"].Byte())
	require.Equal(t, uint16(257), qr.ServerRules.Rules["rule 2"].Uint16())
	require.Equal(t, uint32(16777217), qr.ServerRules.Rules["rule 3"].Uint32())
	require.Equal(t, uint64(72057594037927937), qr.ServerRules.Rules["rule 4"].Uint64())
	require.Equal(t, "string", qr.ServerRules.Rules["rule 5"].String())
}

func testQueryPlayerInfoSinglePacket(t *testing.T, c *queryer) {
	r, err := c.Query()
	require.NoError(t, err, "query request should not have failed")
	qr := r.(*QueryResponse)

	require.Equal(t, uint32(256), c.challengeID, "expected correct challenge id")

	require.NotNil(t, qr, "expected query response")
	require.NotNil(t, qr.PlayerInfo, "expected player info")
	require.Len(t, qr.PlayerInfo.Players, 2)

	require.Len(t, qr.PlayerInfo.Players[0], 5)
	require.Equal(t, byte(128), qr.PlayerInfo.Players[0]["field1"].Byte())
	require.Equal(t, uint16(257), qr.PlayerInfo.Players[0]["field2"].Uint16())
	require.Equal(t, uint32(16777217), qr.PlayerInfo.Players[0]["field3"].Uint32())
	require.Equal(t, uint64(72057594037927937), qr.PlayerInfo.Players[0]["field4"].Uint64())
	require.Equal(t, "string", qr.PlayerInfo.Players[0]["field5"].String())

	require.Len(t, qr.PlayerInfo.Players[1], 5)
	require.Equal(t, byte(129), qr.PlayerInfo.Players[1]["field1"].Byte())
	require.Equal(t, uint16(258), qr.PlayerInfo.Players[1]["field2"].Uint16())
	require.Equal(t, uint32(16777218), qr.PlayerInfo.Players[1]["field3"].Uint32())
	require.Equal(t, uint64(72057594037927938), qr.PlayerInfo.Players[1]["field4"].Uint64())
	require.Equal(t, "STRING", qr.PlayerInfo.Players[1]["field5"].String())
}

func testQueryTeamInfoSinglePacket(t *testing.T, c *queryer) {
	r, err := c.Query()
	require.NoError(t, err, "query request should not have failed")
	qr := r.(*QueryResponse)

	require.Equal(t, uint32(256), c.challengeID, "expected correct challenge id")

	require.NotNil(t, qr, "expected query response")
	require.NotNil(t, qr.TeamInfo, "expected Team info")
	require.Len(t, qr.TeamInfo.Teams, 2)

	require.Len(t, qr.TeamInfo.Teams[0], 5)
	require.Equal(t, byte(128), qr.TeamInfo.Teams[0]["field1"].Byte())
	require.Equal(t, uint16(257), qr.TeamInfo.Teams[0]["field2"].Uint16())
	require.Equal(t, uint32(16777217), qr.TeamInfo.Teams[0]["field3"].Uint32())
	require.Equal(t, uint64(72057594037927937), qr.TeamInfo.Teams[0]["field4"].Uint64())
	require.Equal(t, "string", qr.TeamInfo.Teams[0]["field5"].String())

	require.Len(t, qr.TeamInfo.Teams[1], 5)
	require.Equal(t, byte(129), qr.TeamInfo.Teams[1]["field1"].Byte())
	require.Equal(t, uint16(258), qr.TeamInfo.Teams[1]["field2"].Uint16())
	require.Equal(t, uint32(16777218), qr.TeamInfo.Teams[1]["field3"].Uint32())
	require.Equal(t, uint64(72057594037927938), qr.TeamInfo.Teams[1]["field4"].Uint64())
	require.Equal(t, "STRING", qr.TeamInfo.Teams[1]["field5"].String())
}
