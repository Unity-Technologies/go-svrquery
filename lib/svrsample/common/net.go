package common

import (
	"fmt"
	"net"
	"time"
)

var (
	_ Transport = (*UDPTransport)(nil)
	_ Transport = (*HTTPTransport)(nil)
)

type Transport interface {
	Setup() error
	Addr() string
	Read([]byte) error
	Write([]byte) error
}

type UDPTransport struct {
	address    string
	udpAddress *net.UDPAddr
	conn       *net.UDPConn
}

func NewUDPTransport(address string) UDPTransport {
	return UDPTransport{address: address}
}

func (u UDPTransport) Setup() error {
	addr, err := net.ResolveUDPAddr("udp4", u.address)
	if err != nil {
		return fmt.Errorf("resolved udp: %w", err)
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return fmt.Errorf("listen udp: %w", err)
	}
	u.conn = conn

	return nil
}

func (u UDPTransport) Addr() string {
	return u.udpAddress.String()
}

func (u UDPTransport) Read(buf []byte) error {
	_, udpAddr, err := u.conn.ReadFromUDP(buf)
	if err != nil {
		return fmt.Errorf("read udp: %w", err)
	}
	u.udpAddress = udpAddr
	return nil
}

func (u UDPTransport) Write(resp []byte) error {
	if err := u.conn.SetWriteDeadline(time.Now().Add(1 * time.Second)); err != nil {
		return fmt.Errorf("set write deadline: %w", err)
	}

	if _, err := u.conn.WriteTo(resp, u.udpAddress); err != nil {
		return fmt.Errorf("write udp: %w", err)
	}
	return nil
}

type HTTPTransport struct {
	address string
}

func NewHTTPTransport(address string) UDPTransport {
	return HTTPTransport{address: address}
}

func (h HTTPTransport) Setup() error {
	//TODO implement me
	panic("implement me")
}

func (h HTTPTransport) Addr() string {
	//TODO implement me
	panic("implement me")
}

func (h HTTPTransport) Read(buf []byte) error {
	//TODO implement me
	panic("implement me")
}

func (h HTTPTransport) Write(resp []byte) error {
	//TODO implement me
	panic("implement me")
}
