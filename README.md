# go-svrquery [![Go Report Card](https://goreportcard.com/badge/github.com/multiplay/go-svrquery)](https://goreportcard.com/report/github.com/multiplay/go-svrquery) [![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/multiplay/go-svrquery/blob/master/LICENSE) [![GoDoc](https://godoc.org/github.com/multiplay/go-svrquery?status.svg)](https://godoc.org/github.com/multiplay/go-svrquery) [![Build Status](https://travis-ci.org/multiplay/go-svrquery.svg?branch=master)](https://travis-ci.org/multiplay/go-svrquery)

go-svrquery is a [Go](http://golang.org/) client for talking to game servers using various query protocols.

Features
--------
* Support for various game server query protocol's including:
** SQP, TF2E, Prometheus

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

### Prometheus

To scape metrics exposed in Prometheus/OpenMetrics format, use the protocol `prom`, and the URL of the HTTP metrics 
endpoint. For example:

```go
client, err := svrquery.NewClient("prom", "http://127.0.0.1:9100/metrics")
```

The library expects the following metrics to be available:

```text
# HELP current_players Number of players currently connected to the server.
# TYPE current_players gauge
current_players 1
# HELP max_players Maximum number of players that can connect to the server.
# TYPE max_players gauge
max_players 2
# HELP server_info Server status info.
# TYPE server_info gauge
server_info{game_type="Game Type",map_name="Map",port="1000",server_name="Name"} 1
```

CLI
-------------
A cli is available in github releases and also at https://github.com/multiplay/go-svrquery/tree/master/cmd/cli

This enables you make queries to servers using the specified protocol, and returns the response in pretty json.

### Client

```
./go-svrquery -addr localhost:12121 -proto sqp
{
        "version": 1,
        "address": "localhost:12121",
        "server_info": {
                "current_players": 1,
                "max_players": 2,
                "server_name": "Name",
                "game_type": "Game Type",
                "build_id": "",
                "map": "Map",
                "port": 1000
        }
}
```

### Example Server

This tool also provides the ability to start a very basic sample server using a given protocol.

Currently, only `sqp` and `prom` are supported

```
./go-svrquery -server :12121 -proto sqp
Starting sample server using protocol sqp on :12121
```

Documentation
-------------
- [GoDoc API Reference](http://godoc.org/github.com/multiplay/go-svrquery).

License
-------
go-svrquery is available under the [BSD 2-Clause License](https://opensource.org/licenses/BSD-2-Clause).
```
