[![Build Status](https://travis-ci.org/aopoltorzhicky/go_kraken.svg?branch=master)](https://travis-ci.org/aopoltorzhicky/go_kraken)
[![codecov](https://codecov.io/gh/aopoltorzhicky/go_kraken/branch/master/graph/badge.svg)](https://codecov.io/gh/aopoltorzhicky/go_kraken)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/13668a45df3841b2803cb167beca5032)](https://www.codacy.com/app/aopoltorzhicky/go_kraken?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=aopoltorzhicky/go_kraken&amp;utm_campaign=Badge_Grade)

# Kraken Go
Go library for Kraken Websocket and REST API

## Installation Websocket package

```bash
go get github.com/aopoltorzhicky/go_kraken/websocket
```

## Installation REST API package

Now only Public API realized. Private API is under developing.

```bash
go get github.com/aopoltorzhicky/go_kraken/rest
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


