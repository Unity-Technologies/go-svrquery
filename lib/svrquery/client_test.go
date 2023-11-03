package svrquery

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	cases := []struct {
		name     string
		protocol string
		addr     string
		err      bool
	}{
		{
			name:     "tf2e",
			protocol: "tf2e",
		},
		{
			name:     "prometheus",
			protocol: "prom",
		},
		{
			name:     "invalid-protocol",
			protocol: "my-protocol",
			err:      true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := NewClient(tc.protocol, tc.addr, WithKey("test"), WithTimeout(time.Second))
			if tc.err {
				require.Error(t, err)
				require.Nil(t, c)
			} else {
				require.NoError(t, err)
				require.NotNil(t, c)
			}
		})
	}
}

func TestRead(t *testing.T) {
	t.Fatal("unimplemented")
}

func TestQuery(t *testing.T) {
	addr := os.Getenv("TEST_QUERY_ADDR")
	if addr == "" {
		t.Skip("env TEST_QUERY_ADDR not set")
	}

	proto := os.Getenv("TEST_QUERY_PROTO")
	if proto == "" {
		t.Skip("env TEST_QUERY_PROTO not set")
	}

	c, err := NewClient(proto, addr)
	require.NoError(t, err)
	for i := 0; i < 5; i++ {
		r, err := c.Query()
		require.NoError(t, err)
		fmt.Printf("%#v\n", r)
	}
}
