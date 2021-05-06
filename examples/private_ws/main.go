package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	ws "github.com/aopoltorzhicky/go_kraken/websocket"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	kraken := ws.NewKraken(ws.AuthSandboxBaseURL)
	if err := kraken.Connect(); err != nil {
		log.Fatalf("Error connecting to web socket: %s", err.Error())
	}

	if err := kraken.Authenticate(os.Getenv("KRAKEN_API_KEY"), os.Getenv("KRAKEN_SECRET")); err != nil {
		log.Fatalf("Authenticate error: %s", err.Error())
	}

	if err := kraken.SubscribeOwnTrades(); err != nil {
		log.Fatalf("SubscribeOwnTrades error: %s", err.Error())
	}

	for {
		select {
		case <-signals:
			log.Warn("Stopping...")
			if err := kraken.Close(); err != nil {
				log.Fatal(err)
			}
			return
		case update := <-kraken.Listen():
			switch data := update.Data.(type) {
			case ws.OwnTradesUpdate:
				for i := range data {
					for tradeID, trade := range data[i] {
						log.Printf("Trade %s: %s", tradeID, trade.Type)
					}
				}
			case ws.OpenOrdersUpdate:
				for i := range data {
					for orderID, order := range data[i] {
						log.Printf("Order %s: %#v", orderID, order.Descr)
					}
				}
			default:
			}
		}
	}
}
