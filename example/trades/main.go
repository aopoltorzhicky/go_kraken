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

	// subscribe to BTCUSD trades
	ctx, cxl2 := context.WithTimeout(context.Background(), time.Second*5)
	defer cxl2()
	err = c.SubscribeTrades(ctx, []string{"BTC/USD"})
	if err != nil {
		log.Fatal(err)
	}

	for obj := range c.Listen() {
		switch obj.(type) {
		case error:
			log.Printf("channel closed: %s", obj)
		default:
		}
		for _, trade := range obj.([]ws.TradeUpdate) {
			log.Print("----------------")
			log.Printf("Price: %f", trade.Price)
			log.Printf("Volume: %f", trade.Volume)
			log.Printf("Time: %s", trade.Time.String())
			log.Printf("Pair: %s", trade.Pair)
			log.Printf("Order type: %s", trade.OrderType)
			log.Printf("Side: %s", trade.Side)
			log.Printf("Misc: %s", trade.Misc)
		}
	}
}
