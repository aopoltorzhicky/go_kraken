# Kraken Go
Go library for Kraken Websocket and REST API.

**ATTENTION!** Version 0.1.1 of WebSocket API package is available now! It's not compatible with previous versions. **Please check your code after package update!**

## Installation Websocket package

```bash
go get github.com/aopoltorzhicky/go_kraken/websocket
```

## Installation REST API package

Now only Public API realized

```bash
go get github.com/aopoltorzhicky/go_kraken/rest
```

## Usage

To learn how you can use the package read [examples](examples/).


### Websocket API

For quick start read the one below:

```go
package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	kraken := ws.NewKraken(ws.ProdBaseURL)
	if err := kraken.Connect(); err != nil {
		log.Fatalf("Error connecting to web socket: %s", err.Error())
	}

	// subscribe to BTCUSD`s ticker
	if err := kraken.SubscribeTicker([]string{ws.BTCUSD}); err != nil {
		log.Fatalf("SubscribeTicker error: %s", err.Error())
	}

	for {
		select {
		case <-signals:
			log.Warn("Stopping...")
			if err := kraken.Close(); err != nil {
				log.Fatal(err)
			}
			return
		case update := <-kraken.Listen():
			switch data := update.Data.(type) {
			case ws.TickerUpdate:
				log.Printf("----Ticker of %s----", update.Pair)
				log.Printf("Ask: %s with %s", data.Ask.Price.String(), data.Ask.Volume.String())
				log.Printf("Bid: %s with %s", data.Bid.Price.String(), data.Bid.Volume.String())
				log.Printf("Open today: %s | Open last 24 hours: %s", data.Open.Today.String(), data.Open.Last24.String())
			default:
			}
		}
	}
}
```

Some options is available for `Kraken` object:
```go
kraken := ws.NewKraken(
	ws.ProdBaseURL,
	ws.WithHeartbeatTimeout(10*time.Second), // set interval ping message sending. Should be less than read timeout. Default: 10s.
	ws.WithLogLevel(log.TraceLevel), // set logging level. Default: info.
	ws.WithReadTimeout(15*time.Second), // set read timeout. Default: 15s.
	ws.WithReconnectTimeout(5*time.Second),  // set interval of reconnecting after disconnect. Default: 5s.
)
```

For private Webscoket API usage:
```go
package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Create `Kraken` object
	kraken := ws.NewKraken(ws.AuthSandboxBaseURL)

	// Connect to server
	if err := kraken.Connect(); err != nil {
		log.Fatalf("Error connecting to web socket: %s", err.Error())
	}

	// Authenticate with your private keys
	if err := kraken.Authenticate(os.Getenv("KRAKEN_API_KEY"), os.Getenv("KRAKEN_SECRET")); err != nil {
		log.Fatalf("Authenticate error: %s", err.Error())
	}

	// Subscribe to channels or send commands
	if err := kraken.SubscribeOwnTrades(); err != nil {
		log.Fatalf("SubscribeOwnTrades error: %s", err.Error())
	}

	for {
		select {
		case <-signals:
			log.Warn("Stopping...")
			if err := kraken.Close(); err != nil {
				log.Fatal(err)
			}
			return
		case update := <-kraken.Listen():
			switch data := update.Data.(type) {
			case ws.OwnTradesUpdate:
				for i := range data {
					for tradeID, trade := range data[i] {
						log.Printf("Trade %s: %s", tradeID, trade.Type)
					}
				}
			case ws.OpenOrdersUpdate:
				for i := range data {
					for orderID, order := range data[i] {
						log.Printf("Order %s: %#v", orderID, order.Descr)
					}
				}
			default:
			}
		}
	}
}
```

### REST API

To learn how to use REST API read example below:

```go
package main

import (
	"log"

	"github.com/aopoltorzhicky/go_kraken/rest"
)

func main() {
	api := rest.New("", "")
	spread, err := api.GetSpread("ADAETH", 0)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(spread)

	t, err := api.Time()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(t)
}

```


