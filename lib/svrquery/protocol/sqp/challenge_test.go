package sqp

import (
	"testing"

	"github.com/multiplay/go-svrquery/lib/svrquery/clienttest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestChallenge(t *testing.T) {
	cases := []struct {
		name string
		f    func(t *testing.T, c *queryer)
	}{
		{
			name: "success",
			f: func(t *testing.T, c *queryer) {
				require.NoError(t, c.Challenge(), "challenge request failed")
				require.Equal(t, uint32(256), c.challengeID)
			},
		},
		{
			name: "invalid",
			f: func(t *testing.T, c *queryer) {
				require.Error(t, c.Challenge(), "expected challenge request failure")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := clienttest.LoadData(t, testDir, "challenge_"+tc.name+"_request")
			resp := clienttest.LoadData(t, testDir, "challenge_"+tc.name+"_response")

			m, c := newClient(ServerInfo)
			m.On("Write", req).Return(len(req), nil)
			m.On("Read", mock.AnythingOfType("[]uint8")).Return(resp, nil)
			tc.f(t, c)
		})
	}
}
