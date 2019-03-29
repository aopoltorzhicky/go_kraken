[![Build Status](https://travis-ci.org/aopoltorzhicky/go_kraken.svg?branch=master)](https://travis-ci.org/aopoltorzhicky/go_kraken)
[![codecov](https://codecov.io/gh/aopoltorzhicky/go_kraken/branch/master/graph/badge.svg)](https://codecov.io/gh/aopoltorzhicky/go_kraken)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/80b48b4a7c0a462aba03839821126f0c)](https://www.codacy.com/app/aopoltorzhicky/go_kraken?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=aopoltorzhicky/go_kraken&amp;utm_campaign=Badge_Grade)

# Kraken WebSocket
Go library for Kraken Websocket

## Installation

```bash
go get github.com/aopoltorzhicky/go_kraken/websocket
```

## Usage

To learn how you can use the package read [examples](examples/).

For quick start read the one below:

```go
package main

import (
	"log"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
)

func main() {
	c := ws.New(false)
	err := c.Connect()
	if err != nil {
		log.Fatal("Error connecting to web socket : ", err)
	}

	err = c.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// subscribe to ADABTC, XTZBTC, XLMBTC spread
	err = c.SubscribeSpread([]string{ws.ADABTC, ws.XTZBTC, ws.XLMBTC})
	if err != nil {
		log.Fatal(err)
	}

	for obj := range c.Listen() {
		switch obj.(type) {
		case error:
			log.Printf("channel closed: %s", obj)
		default:
		}
		log.Printf("MSG RECV: %#v", obj)
	}
}

```


