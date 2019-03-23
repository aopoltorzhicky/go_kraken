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

	pairs := []string{ws.BTCUSD}
	// subscribe to BTCUSD trades
	ctx, cxl := context.WithTimeout(context.Background(), time.Second*5)
	defer cxl()
	err = c.SubscribeTrades(ctx, pairs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		time.Sleep(time.Second * 2)
		log.Print("Unsubsribing...")
		ctx2, cxl2 := context.WithTimeout(context.Background(), time.Second*5)
		defer cxl2()
		err = c.Unsubscribe(ctx2, ws.ChanTrades, pairs)
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
