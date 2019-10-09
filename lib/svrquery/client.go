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

// Client provides the ability to query a server.
type Client struct {
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
func WithKey(key string) func(*Client) error {
	return func(c *Client) error {
		c.key = key
		return nil
	}
}

// WithTimeout sets the read and write timeout for the client.
func WithTimeout(t time.Duration) func(*Client) error {
	return func(c *Client) error {
		c.timeout = t
		return nil
	}
}

// NewClient creates a new client that talks to addr.
func NewClient(proto, addr string, options ...func(*Client) error) (*Client, error) {
	f, err := protocol.Get(proto)
	if err != nil {
		return nil, err
	}

	c := &Client{
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
func (c *Client) Write(b []byte) (int, error) {
	if err := c.c.SetWriteDeadline(time.Now().Add(c.timeout)); err != nil {
		return 0, err
	}

	return c.c.Write(b)
}

// Read implements io.Reader.
func (c *Client) Read(b []byte) (int, error) {
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
func (c *Client) Close() error {
	return c.c.Close()
}

// Key implements common.Client.
func (c *Client) Key() string {
	return c.key
}

// Protocol returns the protocol of the client.
func (c *Client) Protocol() string {
	return c.protocol
}
