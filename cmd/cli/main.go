package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/multiplay/go-svrquery/lib/svrquery"
	"github.com/multiplay/go-svrquery/lib/svrsample"
	"github.com/multiplay/go-svrquery/lib/svrsample/common"
	"log"
	"os"
)

func main() {
	clientAddr := flag.String("addr", "", "Address to connect to e.g. 127.0.0.1:12345")
	proto := flag.String("proto", "", "Protocol e.g. sqp, tf2e, tf2e-v7, tf2e-v8")
	key := flag.String("key", "", "Key to use to authenticate")
	file := flag.String("file", "", "Bulk file to execute to get basic server information")
	serverAddr := flag.String("server", "", "Address to start server e.g. 127.0.0.1:12121, :23232")
	flag.Parse()

	l := log.New(os.Stderr, "", 0)

	if *file != "" {
		// Use bulk file mode
		if err := queryBulk(*file); err != nil {
			l.Fatal(err)
		}
		return
	}

	if *serverAddr != "" && *clientAddr != "" {
		bail(l, "Cannot run both a server and a client. Specify either -addr OR -server flags")
	}

	switch {
	case *serverAddr != "":
		if *proto == "" {
			bail(l, "No protocol provided in client mode")
		}
		serverMode(l, *proto, *serverAddr)
	case *clientAddr != "":
		if *proto == "" {
			bail(l, "Protocol required in server mode")
		}
		queryMode(l, *proto, *clientAddr, *key)
	default:
		bail(l, "Please supply some options")
	}
}

func queryMode(l *log.Logger, proto, address, key string) {
	if err := query(proto, address, key); err != nil {
		l.Fatal(err)
	}
}

func query(proto, address, key string) error {
	options := make([]svrquery.Option, 0)
	if key != "" {
		options = append(options, svrquery.WithKey(key))
	}

	c, err := svrquery.NewClient(proto, address, options...)
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

	transport, err := svrsample.GetTransport(proto, address)
	if err != nil {
		return err
	}
	err = transport.Setup()
	if err != nil {
		return err
	}

	responder, err := svrsample.GetResponder(proto, common.QueryState{
		CurrentPlayers: 1,
		MaxPlayers:     2,
		ServerName:     "Name",
		GameType:       "Game Type",
		Map:            "Map",
		Port:           1000,
	})
	if err != nil {
		return err
	}

	for {
		buf := make([]byte, 16)
		err = transport.Read(buf)
		if err != nil {
			l.Println("read", err)
			continue
		}

		resp, err := responder.Respond(transport.Addr(), buf)
		if err != nil {
			l.Println("responding to query", err)
			continue
		}

		if err = transport.Write(resp); err != nil {
			l.Println("writing response")
		}
	}

}

func bail(l *log.Logger, msg string) {
	l.Println(msg)
	flag.PrintDefaults()
	os.Exit(1)
}
