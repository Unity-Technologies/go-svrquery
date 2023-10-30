package svrquery

import (
	"net"
	"time"

	"github.com/multiplay/go-svrquery/lib/svrquery/protocol"
	// Register all known protocols
	_ "github.com/multiplay/go-svrquery/lib/svrquery/protocol/all"
)

var (
	// DefaultTimeout is the default read and write timeout.
	DefaultTimeout = time.Millisecond * 1000

	// DefaultNetwork is the default network for a new client.
	DefaultNetwork = "udp"
)

var (
	_ protocol.Client = (*UDPClient)(nil)
	_ protocol.Client = (*HTTPClient)(nil)
)

// Option represents a UDPClient option.
type Option func(*UDPClient) error

// UDPClient provides the ability to query a server.
type UDPClient struct {
	protocol string
	network  string
	addr     string
	ua       *net.UDPAddr
	key      string
	timeout  time.Duration
	c        *net.UDPConn
	protocol.Queryer
}

// WithKey sets the key used for request by for the client.
func WithKey(key string) Option {
	return func(c *UDPClient) error {
		c.key = key
		return nil
	}
}

// WithTimeout sets the read and write timeout for the client.
func WithTimeout(t time.Duration) Option {
	return func(c *UDPClient) error {
		c.timeout = t
		return nil
	}
}

// NewUDPClient creates a new client that talks to addr.
func NewUDPClient(proto, addr string, options ...Option) (*UDPClient, error) {
	f, err := protocol.Get(proto)
	if err != nil {
		return nil, err
	}

	c := &UDPClient{
		protocol: proto,
		addr:     addr,
		network:  DefaultNetwork,
		timeout:  DefaultTimeout,
	}
	c.Queryer = f(c)

	for _, o := range options {
		if err := o(c); err != nil {
			return nil, err
		}
	}

	if c.ua, err = net.ResolveUDPAddr(c.network, addr); err != nil {
		return nil, err
	}

	if c.c, err = net.DialUDP(c.network, nil, c.ua); err != nil {
		return nil, err
	}

	return c, nil
}

// Write implements io.Writer.
func (c *UDPClient) Write(b []byte) (int, error) {
	if err := c.c.SetWriteDeadline(time.Now().Add(c.timeout)); err != nil {
		return 0, err
	}

	return c.c.Write(b)
}

// Read implements io.Reader.
func (c *UDPClient) Read(b []byte) (int, error) {
	if err := c.c.SetReadDeadline(time.Now().Add(c.timeout)); err != nil {
		return 0, err
	}

	for {
		n, addr, err := c.c.ReadFromUDP(b)
		if err != nil {
			return 0, err
		} else if addr.String() == c.ua.String() { // We use String as IP's can be different byte but the same value.
			return n, nil
		}
		// Packet from unexpected source just ignore.
	}
}

// Close implements io.Closer.
func (c *UDPClient) Close() error {
	return c.c.Close()
}

// Key implements protocol.UDPClient.
func (c *UDPClient) Key() string {
	return c.key
}

// Address implements protocol.UDPClient.
func (c *UDPClient) Address() string {
	return c.addr
}

// Protocol returns the protocol of the client.
func (c *UDPClient) Protocol() string {
	return c.protocol
}

type HTTPClient struct {
	protocol.Queryer
	address string
}

func NewHTTPClient(proto, address string) (*HTTPClient, error) {
	queryerCreator, err := protocol.Get(proto)
	if err != nil {
		return nil, err
	}
	client := &HTTPClient{address: address}
	client.Queryer = queryerCreator(client)

	return client, nil
}

func (c *HTTPClient) Read(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *HTTPClient) Write(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (c *HTTPClient) Query() (protocol.Responser, error) {
	return c.Queryer.Query()
}

func (c *HTTPClient) Key() string {
	//TODO implement me
	panic("implement me")
}

func (c *HTTPClient) Address() string {
	return c.address
}

// Close implements io.Closer.
func (c *HTTPClient) Close() error {
	// no-op
	return nil
}
