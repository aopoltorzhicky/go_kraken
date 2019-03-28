[![Build Status](https://travis-ci.org/aopoltorzhicky/go_kraken.svg?branch=master)](https://travis-ci.org/aopoltorzhicky/go_kraken)

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

	ws "scripts/kraken_ws/websocket"
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

	// subscribe to BTCUSD, XLMUSD, ADACAD spread
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


