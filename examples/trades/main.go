package main

import (
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

	pairs := []string{ws.BTCUSD}

	// subscribe to BTCUSD trades
	err = c.SubscribeTrades(pairs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		time.Sleep(time.Second * 2)
		log.Print("Unsubsribing...")
		err = c.Unsubscribe(ws.ChanTrades, pairs)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Success!")
		c.Close()
	}()

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
