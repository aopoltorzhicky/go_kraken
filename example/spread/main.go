package main

import (
	"context"
	"log"
	"time"

	ws "scripts/kraken_ws/websocket"
)

func main() {
	c := ws.New(false)
	err := c.Connect()
	if err != nil {
		log.Fatal("Error connecting to web socket : ", err)
	}

	// subscribe to BTCUSD, XLMUSD, ADACAD spread
	ctx, cxl2 := context.WithTimeout(context.Background(), time.Second*5)
	defer cxl2()
	err = c.SubscribeSpread(ctx, []string{ws.ADABTC, ws.XTZBTC, ws.XLMBTC})
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
