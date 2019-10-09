package main

import (
	"log"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
)

func main() {
	c := ws.New()
	err := c.Connect()
	if err != nil {
		log.Fatal("Error connecting to web socket : ", err)
	}

	// subscribe to BTCUSD ticker
	err = c.SubscribeTicker([]string{ws.BTCUSD})
	if err != nil {
		log.Fatal(err)
	}

	for obj := range c.Listen() {
		switch obj.(type) {
		case error:
			log.Printf("channel closed: %s", obj)
		case ws.DataUpdate:
			data := obj.(ws.DataUpdate)
			ticker := data.Data.(ws.TickerUpdate)
			log.Print("-------------")
			log.Printf("Ask: %f with %f", ticker.Ask.Price, ticker.Ask.Volume)
			log.Printf("Bid: %f with %f", ticker.Bid.Price, ticker.Bid.Volume)
			log.Printf("Open today: %f | Open last 24 hours: %f", ticker.Open.Today, ticker.Open.Last24)
		default:
		}

	}
}
