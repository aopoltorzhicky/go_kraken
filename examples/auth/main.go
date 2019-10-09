package main

import (
	"log"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
)

func main() {
	c := ws.NewAuth("", "")

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
