package clienttest

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockClient struct {
	mock.Mock
}

func (mc *MockClient) Address() string {
	args := mc.Called()
	return args.String(0)
}

func (mc *MockClient) Args() map[string]interface{} {
	args := mc.Called()
	return args.Get(0).(map[string]interface{})
}

func (mc *MockClient) Write(b []byte) (int, error) {
	args := mc.Called(b)
	return args.Int(0), args.Error(1)
}

func (mc *MockClient) Read(b []byte) (int, error) {
	args := mc.Called(b)
	d := args.Get(0).([]byte)
	copy(b, d)
	return len(d), args.Error(1)
}

func (mc *MockClient) Close() error {
	args := mc.Called()
	return args.Error(0)
}

func (mc *MockClient) Key() string {
	args := mc.Called()
	return args.String(0)
}

func LoadData(t *testing.T, fileParts ...string) []byte {
	d, err := ioutil.ReadFile(filepath.Join(fileParts...))
	require.NoError(t, err)
	return d
}

func LoadMultiData(t *testing.T, files int, fileParts ...string) [][]byte {
	pkts := make([][]byte, files)
	file := fileParts[len(fileParts)-1]
	for i := range pkts {
		fileParts[len(fileParts)-1] = fmt.Sprintf("%s_%03d", file, i)
		d, err := ioutil.ReadFile(filepath.Join(fileParts...))
		require.NoError(t, err)
		pkts[i] = d
		require.NoError(t, err)
	}
	return pkts
}
