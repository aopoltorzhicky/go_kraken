package main

import (
	"log"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
)

func main() {
	c := ws.NewAuth("l+GpEJ7HT+8qOQJnNGl4rFTwcWYuYb8I1ua3n1MrWthsD67/aui7UQ7p", "LQ7LzKfLnRU8iKZrhKtIF74nd5tibb0tlCQr3JyoWS7RjK8bhKCO5xEV/FZiQTggqTiDQOCEQvNIrwZ2S+GnQA==", true)

	if err := c.Connect(); err != nil {
		log.Fatal("Error connecting to web socket : ", err)
	}

	if err := c.SubscribeOwnTrades(); err != nil {
		log.Fatal("Error connecting to web socket : ", err)
	}

	if err := c.SubscribeOpenOrders(); err != nil {
		log.Fatal("Error connecting to web socket : ", err)
	}

	for obj := range c.Listen() {
		switch obj.(type) {
		case error:
			log.Printf("channel closed: %s", obj)
		case ws.DataUpdate:
			data := obj.(ws.DataUpdate)
			log.Printf("MSG RECV: %#v", data)
		}
	}
}
