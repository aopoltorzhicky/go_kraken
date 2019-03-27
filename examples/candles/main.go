package main

import (
	"log"

	ws "githib.com/aopoltorzhicky/go_kraken/websocket"
)

func main() {
	c := ws.New(false)
	err := c.Connect()
	if err != nil {
		log.Fatal("Error connecting to web socket : ", err)
	}

	// subscribe to BTCUSD, XLMUSD, ADACAD candles
	err = c.SubscribeCandles([]string{ws.XLMBTC, ws.BTCEUR, ws.QTUMCAD}, ws.Interal10080)
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
