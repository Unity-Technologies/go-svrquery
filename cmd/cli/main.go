package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/multiplay/go-svrquery/lib/svrquery"
)

func main() {
	address := flag.String("addr", "", "Address e.g. 127.0.0.1:12345")
	proto := flag.String("proto", "", "Protocol e.g. sqp, tf2e")
	flag.Parse()

	l := log.New(os.Stderr, "", 0)

	if *address == "" {
		l.Println("No address provided")
		flag.PrintDefaults()
		os.Exit(1)

	}

	if *proto == "" {
		l.Println("No protocol provided")
		flag.PrintDefaults()
		os.Exit(1)
	}

	c, err := svrquery.NewClient(*proto, *address)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	r, err := c.Query()
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s \n", b)
}
