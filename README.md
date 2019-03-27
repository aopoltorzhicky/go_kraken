# Kraken WS
Go library for Kraken Websocket

## Installation

```bash
go get github.com/aopoltorzhicky/kraken_ws
```

## Usage

To learn how you can use the package read [examples](kraken_ws/blob/master/examples/).

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


