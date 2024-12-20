package svrsample

import (
	"context"
	"errors"
	"fmt"
	"github.com/multiplay/go-svrquery/lib/svrsample/common"
	"github.com/multiplay/go-svrquery/lib/svrsample/protocol/prom"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	_ Transport = (*UDPTransport)(nil)
	_ Transport = (*HTTPTransport)(nil)
)

// Transport is an abstraction of the metrics transport (UDP, HTTP, etc.)
type Transport interface {
	// Start starts the transport and blocks until it is stopped
	Start(context.Context, common.QueryResponder) error
}

type UDPTransport struct {
	address    string
	udpAddress *net.UDPAddr
	conn       *net.UDPConn
}

func NewUDPTransport(address string) UDPTransport {
	return UDPTransport{address: address}
}

func (u UDPTransport) Start(ctx context.Context, responder common.QueryResponder) error {
	// TODO: do something with context

	addr, err := net.ResolveUDPAddr("udp4", u.address)
	if err != nil {
		return fmt.Errorf("resolved udp: %w", err)
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return fmt.Errorf("listen udp: %w", err)
	}
	u.conn = conn

	for {
		buf := make([]byte, 16)
		err := u.read(buf)
		if err != nil {
			log.Println("read", err)
			continue
		}

		resp, err := responder.Respond(u.udpAddress.String(), buf)
		if err != nil {
			log.Println("responding to query", err)
			continue
		}

		if err = u.write(resp); err != nil {
			log.Println("writing response")
		}
	}
}

func (u UDPTransport) read(buf []byte) error {
	_, udpAddr, err := u.conn.ReadFromUDP(buf)
	if err != nil {
		return fmt.Errorf("read udp: %w", err)
	}
	u.udpAddress = udpAddr
	return nil
}

func (u UDPTransport) write(resp []byte) error {
	if err := u.conn.SetWriteDeadline(time.Now().Add(1 * time.Second)); err != nil {
		return fmt.Errorf("set write deadline: %w", err)
	}

	if _, err := u.conn.WriteTo(resp, u.udpAddress); err != nil {
		return fmt.Errorf("write udp: %w", err)
	}
	return nil
}

type HTTPTransport struct {
	address    string
	httpServer *http.Server
}

func NewHTTPTransport(address string) HTTPTransport {
	return HTTPTransport{address: address}
}

func (h HTTPTransport) Start(ctx context.Context, responder common.QueryResponder) error {
	promResponder, ok := responder.(*prom.QueryResponder)
	if !ok {
		return errors.New(fmt.Sprintf("bad responder type, expected prom.QueryResponder but got %T", responder))
	}

	// TODO: do something with context

	listener, err := net.Listen("tcp", h.address)
	if err != nil {
		return fmt.Errorf("listen tcp: %w", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promResponder.HTTPHandler)
	httpServer := &http.Server{Addr: h.address, Handler: mux}
	h.httpServer = httpServer

	if err = h.httpServer.Serve(listener); err != nil {
		return fmt.Errorf("serve http: %w", err)
	}
	return nil
}
