package main

import (
	"context"
	"log"
	"time"

	ws "scripts/kraken_ws/websocket"
)

func main() {
	c := ws.New()
	err := c.Connect()
	if err != nil {
		log.Fatal("Error connecting to web socket : ", err)
	}

	// subscribe to BTCUSD ticker
	ctx, cxl2 := context.WithTimeout(context.Background(), time.Second*5)
	defer cxl2()
	err = c.SubscribeTicker(ctx, []string{"BTC/USD"})
	if err != nil {
		log.Fatal(err)
	}

	for obj := range c.Listen() {
		switch obj.(type) {
		case error:
			log.Printf("channel closed: %s", obj)
		default:
		}
		ticker := obj.(ws.TickerUpdate)
		log.Print("-------------")
		log.Printf("Ask: %f with %f", ticker.Ask.Price, ticker.Ask.Volume)
		log.Printf("Bid: %f with %f", ticker.Bid.Price, ticker.Bid.Volume)
		log.Printf("Open today: %f | Open last 24 hours: %f", ticker.Open.Today, ticker.Open.Last24)
	}
}
