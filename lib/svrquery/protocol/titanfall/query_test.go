package titanfall

import (
	"testing"

	"github.com/multiplay/go-svrquery/lib/svrquery/clienttest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testDir = "testdata"
	testKey = "AABBCCddeeffgghhkkllmmNN"
)

var (
	base = Info{
		Header: Header{
			Prefix:  -1,
			Command: 78,
			Version: 3,
		},
		InstanceInfo: InstanceInfo{
			Retail:         1,
			InstanceType:   2,
			ClientCRC:      4294967295,
			NetProtocol:    526,
			RandomServerID: 0,
		},
		BuildName:  "R5pc_r5launch_N895_CL450114_2019_10_03_04_00_PM",
		Datacenter: "west europe 2",
		GameMode:   "survival",
		BasicInfo: BasicInfo{
			Port:            37015,
			Platform:        "PC",
			PlaylistVersion: "",
			PlaylistNum:     307,
			PlaylistName:    "des_ranked",
			NumClients:      0,
			MaxClients:      60,
			Map:             "mp_rr_desertlands_64k_x_64k",
		},
		PerformanceInfo: PerformanceInfo{},
		MatchState: MatchState{
			MatchStateV2: MatchStateV2{
				Phase:            2,
				MaxRounds:        1,
				RoundsWonIMC:     0,
				RoundsWonMilitia: 0,
				TimeLimit:        1800,
				TimePassed:       0,
				MaxScore:         50,
			},
		},
	}
)

func TestQuery(t *testing.T) {
	keyed := base
	keyed.Version = 5
	keyed.AverageFrameTime = 1.2347187
	keyed.MaxFrameTime = 1.583148
	keyed.AverageUserCommandTime = 0.9734314
	keyed.MaxUserCommandTime = 7.678111

	v7 := base
	v7.Version = 7
	v7.PlatformPlayers = map[string]byte{
		"ps3": 16,
		"pc":  6,
	}
	v7.PerformanceInfo = PerformanceInfo{
		AverageFrameTime:       1,
		MaxFrameTime:           2,
		AverageUserCommandTime: 3,
		MaxUserCommandTime:     4,
	}
	v7.TeamsLeftWithPlayersNum = 6

	cases := []struct {
		name     string
		version  byte
		request  string
		response string
		key      string
		expected Info
	}{
		{
			name:     "v3",
			version:  3,
			request:  "request-v3",
			response: "response-v3",
			expected: base,
		},
		{
			name:     "v7",
			version:  7,
			request:  "request-v7",
			response: "response-v7",
			expected: v7,
		},
		{
			name:     "keyed",
			version:  5,
			request:  "request-key",
			response: "response-key",
			key:      testKey,
			expected: keyed,
		},
		{
			name:     "keyed_upgrades_lower_version",
			version:  3,
			request:  "request-key",
			response: "response-key",
			key:      testKey,
			expected: keyed,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := clienttest.LoadData(t, testDir, tc.request)
			resp := clienttest.LoadData(t, testDir, tc.response)
			m := &clienttest.MockClient{}

			m.On("Write", req).Return(len(req), nil)
			m.On("Read", mock.AnythingOfType("[]uint8")).Return(resp, nil)
			m.On("Key").Return(tc.key)

			p := newQueryer(tc.version)(m)
			i, err := p.Query()
			require.NoError(t, err)
			require.IsType(t, &Info{}, i)
			require.Equal(t, &tc.expected, i)
			m.AssertExpectations(t)
		})
	}
}
