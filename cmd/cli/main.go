package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/multiplay/go-svrquery/lib/svrquery"
	"github.com/multiplay/go-svrquery/lib/svrsample"
)

func main() {
	clientAddr := flag.String("addr", "", "Address e.g. 127.0.0.1:12345")
	proto := flag.String("proto", "", "Protocol e.g. sqp, tf2e, tf2e-v7, tf2e-v8")
	serverAddr := flag.String("server", "", "Address to start server e.g. 127.0.0.1:12121, :23232")
	flag.Parse()

	l := log.New(os.Stderr, "", 0)

	if *serverAddr != "" && *clientAddr != "" {
		bail(l, "Cannot run both a server and a client. Specify either -addr OR -server flags")
	}

	if *serverAddr != "" {
		if *proto == "" {
			bail(l, "No protocol provided")
		}
		serverMode(l, *proto, *serverAddr)
	} else {
		if *proto == "" {
			bail(l, "Protocol required in server mode")
		}
		if *clientAddr == "" {
			bail(l, "No address provided")
		}
		queryMode(l, *proto, *clientAddr)
	}
}

func queryMode(l *log.Logger, proto, address string) {
	if err := query(proto, address); err != nil {
		l.Fatal(err)
	}
}

func query(proto, address string) error {
	c, err := svrquery.NewClient(proto, address)
	if err != nil {
		return err
	}
	defer c.Close()

	r, err := c.Query()
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", b)
	return nil
}

func serverMode(l *log.Logger, proto, serverAddr string) {
	if err := server(l, proto, serverAddr); err != nil {
		l.Fatal(err)
	}
}

func server(l *log.Logger, proto, address string) error {
	l.Printf("Starting sample server using protocol %s on %s", proto, address)
	responder, err := svrsample.GetResponder(proto, svrsample.QueryState{
		CurrentPlayers: 1,
		MaxPlayers:     2,
		ServerName:     "Name",
		GameType:       "Game Type",
		Map:            "Map",
		Port:           1000,
	})

	addr, err := net.ResolveUDPAddr("udp4", address)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		return err
	}

	for {
		buf := make([]byte, 16)
		_, to, err := conn.ReadFromUDP(buf)
		if err != nil {
			l.Println("read from udp", err)
			continue
		}

		resp, err := responder.Respond(to.String(), buf)
		if err != nil {
			l.Println("error responding to query", err)
			continue
		}

		if err = conn.SetWriteDeadline(time.Now().Add(1 * time.Second)); err != nil {
			l.Println("error setting write deadline")
			continue
		}

		if _, err = conn.WriteTo(resp, to); err != nil {
			l.Println("error writing response")
		}
	}

}

func bail(l *log.Logger, msg string) {
	l.Println(msg)
	flag.PrintDefaults()
	os.Exit(1)
}
