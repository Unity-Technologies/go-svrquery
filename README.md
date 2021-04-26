# go-svrquery [![Go Report Card](https://goreportcard.com/badge/github.com/multiplay/go-svrquery)](https://goreportcard.com/report/github.com/multiplay/go-svrquery) [![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/multiplay/go-svrquery/blob/master/LICENSE) [![GoDoc](https://godoc.org/github.com/multiplay/go-svrquery?status.svg)](https://godoc.org/github.com/multiplay/go-svrquery) [![Build Status](https://travis-ci.org/multiplay/go-svrquery.svg?branch=master)](https://travis-ci.org/multiplay/go-svrquery)

go-svrquery is a [Go](http://golang.org/) client for talking to game servers using various query protocols.

Features
--------
* Support for various game server query protocol's including:
** Titanfall
* Supports per protocol custom [netdata](https://github.com/netdata/netdata) graphs as required by [go.d.plugin](https://github.com/netdata/go.d.plugin) which is based on the [go-orchestrator plugin framework](https://github.com/netdata/go-orchestrator).

Installation
------------
```sh
go get -u github.com/multiplay/go-svrquery
```

Examples
--------

Using go-svrquery is simple just create a client and then send commands e.g.
```go
package main

import (
	"log"
	"time"

	"github.com/multiplay/go-svrquery/lib/svrquery"
)

func main() {
	c, err := svrquery.NewClient("tf2e", "192.168.1.102:10011")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	r, err := c.Query()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v\n", r)
}
```

CLI
-------------
A cli is available in https://github.com/multiplay/go-svrquery/tree/master/cmd/cli

This enables you make queries to servers using the specified protocol, and returns the response in pretty json.

Documentation
-------------
- [GoDoc API Reference](http://godoc.org/github.com/multiplay/go-svrquery).

License
-------
go-svrquery is available under the [BSD 2-Clause License](https://opensource.org/licenses/BSD-2-Clause).
```
