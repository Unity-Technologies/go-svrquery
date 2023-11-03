package svrquery

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
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
	_ protocol.Client = (*Client)(nil)
)

// Option represents a Client option.
type Option func(*Client) error

// Client provides the ability to query a server.
type Client struct {
	protocol string
	key      string
	timeout  time.Duration
	protocol.Queryer
	transport
}

// WithKey sets the key used for request by for the client.
func WithKey(key string) Option {
	return func(c *Client) error {
		c.key = key
		return nil
	}
}

// WithTimeout sets the read and write timeout for the client.
func WithTimeout(t time.Duration) Option {
	return func(c *Client) error {
		c.transport.SetTimeout(t)
		return nil
	}
}

// NewClient creates a new client that talks to addr using proto.
func NewClient(proto, addr string, options ...Option) (*Client, error) {
	f, err := protocol.Get(proto)
	if err != nil {
		return nil, err
	}

	c := &Client{
		protocol: proto,
	}
	c.Queryer = f(c)

	switch proto {
	case "prom":
		c.transport = newHTTPTransport(addr)
	default:
		// defaulting to udp
		c.transport = newUDPTransport(addr)
	}

	for _, o := range options {
		if err := o(c); err != nil {
			return nil, err
		}
	}

	if err := c.transport.Setup(); err != nil {
		return nil, fmt.Errorf("setup client transport: %w", err)
	}

	return c, nil
}

func (c *Client) Query() (protocol.Responser, error) {
	return c.Queryer.Query()
}

var (
	_ transport = (*udpTransport)(nil)
	_ transport = (*httpTransport)(nil)
)

type transport interface {
	Setup() error
	Address() string
	SetTimeout(time.Duration)
	io.ReadWriteCloser
}

type udpTransport struct {
	address    string
	timeout    time.Duration
	connection *net.UDPConn
	udpAddress *net.UDPAddr
}

func newUDPTransport(address string) *udpTransport {
	return &udpTransport{address: address}
}

// Address implements transport.Address.
func (u *udpTransport) Address() string {
	return u.address
}

// Setup implements transport.Setup.
func (u *udpTransport) Setup() error {
	udpAddr, err := net.ResolveUDPAddr(DefaultNetwork, u.address)
	if err != nil {
		return err
	}
	u.udpAddress = udpAddr

	if u.connection, err = net.DialUDP(DefaultNetwork, nil, u.udpAddress); err != nil {
		return err
	}
	return nil
}

// Write implements io.Writer.
func (u *udpTransport) Write(b []byte) (int, error) {
	if err := u.connection.SetWriteDeadline(time.Now().Add(u.timeout)); err != nil {
		return 0, err
	}

	return u.connection.Write(b)
}

// Read implements io.Reader.
func (u *udpTransport) Read(b []byte) (int, error) {
	if err := u.connection.SetReadDeadline(time.Now().Add(u.timeout)); err != nil {
		return 0, err
	}

	for {
		n, addr, err := u.connection.ReadFromUDP(b)
		if err != nil {
			return 0, err
		} else if addr.String() == u.udpAddress.String() { // We use String as IP's can be different byte but the same value.
			return n, nil
		}
		// Packet from unexpected source just ignore.
	}
}

// Close implements io.Closer.
func (u *udpTransport) Close() error {
	return u.connection.Close()
}

func (u *udpTransport) SetTimeout(t time.Duration) {
	u.timeout = t
}

// Key implements protocol.Client.
func (c *Client) Key() string {
	return c.key
}

// Protocol returns the protocol of the client.
func (c *Client) Protocol() string {
	return c.protocol
}

type httpTransport struct {
	address    string
	httpClient *http.Client
}

func newHTTPTransport(address string) *httpTransport {
	t := &httpTransport{address: address}
	t.httpClient = &http.Client{}
	return t
}

func (h *httpTransport) Setup() error {
	// no-op
	return nil
}

func (h *httpTransport) Address() string {
	return h.address
}

func (h *httpTransport) Read(b []byte) (int, error) {
	res, err := h.httpClient.Get(h.address)
	if err != nil {
		return 0, fmt.Errorf("http get: %w", err)
	}

	return res.Body.Read(b)
}

func (h *httpTransport) Write(b []byte) (int, error) {
	return 0, errors.New("httpTransport.Write is unused")
}

// Close implements io.Closer.
func (h *httpTransport) Close() error {
	// no-op
	return nil
}

func (h *httpTransport) SetTimeout(t time.Duration) {
	h.httpClient.Timeout = t
}
