package svrquery

import (
	"fmt"
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
	Transport Transport
}

func (c *Client) Read(p []byte) (n int, err error) {
	return c.Transport.Read(p)
}

func (c *Client) Write(p []byte) (n int, err error) {
	return c.Transport.Write(p)
}

func (c *Client) Close() error {
	return c.Transport.Close()
}

func (c *Client) Address() string {
	return c.Transport.Address()
}

type Transport interface {
	Setup() error
	Address() string
	Read(b []byte) (int, error)
	Write(b []byte) (int, error)
	Close() error
	SetTimeout(time.Duration)
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
		c.Transport.SetTimeout(t)
		return nil
	}
}

// NewClient creates a new client that talks to address.
func NewClient(proto, addr string, options ...Option) (*Client, error) {
	f, err := protocol.Get(proto)
	if err != nil {
		return nil, err
	}

	c := &Client{
		protocol: proto,
	}
	c.Queryer = f(c)

	// TODO: move keys & timeout to transport?
	for _, o := range options {
		if err := o(c); err != nil {
			return nil, err
		}
	}

	var t Transport
	switch proto {
	case "sqp":
		t = &udpTransport{address: addr}
	case "prom":
		t = &HTTPTransport{}
	default:
		return nil, fmt.Errorf("protocol %s not supported", proto)
	}

	if err := t.Setup(); err != nil {
		return nil, fmt.Errorf("setup client transport: %w", err)
	}

	c.Transport = t

	return c, nil
}

func (c *Client) Query() (protocol.Responser, error) {
	return c.Queryer.Query()
}

var _ Transport = (*udpTransport)(nil)
var _ Transport = (*HTTPTransport)(nil)

type udpTransport struct {
	address    string
	timeout    time.Duration
	connection *net.UDPConn
	udpAddress *net.UDPAddr
}

// Address implements transport.Address.
func (u *udpTransport) Address() string {
	return u.address
}

// Setup implements transport.Setup.
func (u *udpTransport) Setup() error {
	udpNet := "udp"
	udpAddr, err := net.ResolveUDPAddr(udpNet, u.address)
	if err != nil {
		return err
	}
	u.udpAddress = udpAddr

	if u.connection, err = net.DialUDP(udpNet, nil, u.udpAddress); err != nil {
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

type HTTPTransport struct {
	address    string
	HttpClient *http.Client
}

func (h HTTPTransport) Setup() error {
	h.HttpClient = &http.Client{}
	return nil
}

func (h HTTPTransport) Address() string {
	return h.address
}

func (h HTTPTransport) Read(b []byte) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (h HTTPTransport) Write(b []byte) (int, error) {
	//TODO implement me
	panic("implement me")
}

// Close implements io.Closer.
func (h HTTPTransport) Close() error {
	// no-op
	return nil
}

func (h HTTPTransport) SetTimeout(t time.Duration) {
	h.HttpClient.Timeout = t
}
