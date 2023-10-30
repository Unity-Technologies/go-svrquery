package main

import (
	"testing"

	"github.com/multiplay/go-svrquery/lib/svrquery"
	"github.com/stretchr/testify/require"
)

func TestParseEntry(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		expQuery   string
		expAddress string
		expErr     error
	}{
		{
			name:       "ok",
			input:      "sqp 1.2.3.4:1234",
			expQuery:   "sqp",
			expAddress: "1.2.3.4:1234",
		},
		{
			name:   "empty line",
			input:  "",
			expErr: errNoItem,
		},
		{
			name:   "invalid entry",
			input:  "sqp 1.2.3.4:1234 extra",
			expErr: errEntryInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, addr, err := parseEntry(tc.input)
			if err != nil {
				require.ErrorIs(t, err, tc.expErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expQuery, query)
			require.Equal(t, tc.expAddress, addr)
		})
	}
}

func TestCreateClient(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expQuery string
		expKey   string
		expErr   error
	}{
		{
			name:     "ok",
			query:    "tf2e",
			expQuery: "tf2e",
		},
		{
			name:     "with_key",
			query:    "tf2e,key=val",
			expKey:   "val",
			expQuery: "tf2e",
		},
		{
			name:     "with_unsupported_other",
			query:    "tf2e,other=val",
			expQuery: "tf2e",
		},
		{
			name:     "invalid entry",
			query:    "tf2e",
			expQuery: "tf2e",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			baseQuery, options, err := parseOptions(tc.query)
			if err != nil {
				require.ErrorIs(t, err, tc.expErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expQuery, baseQuery)

			// Validate key setting
			if tc.expKey != "" {
				require.Len(t, options, 1)
				c := svrquery.UDPClient{}
				require.NoError(t, options[0](&c))
				require.Equal(t, tc.expKey, c.Key())
			}
			require.NotNil(t, options)
		})
	}
}
