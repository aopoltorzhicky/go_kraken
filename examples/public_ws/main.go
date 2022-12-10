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

	kraken := ws.NewKraken(ws.ProdBaseURL)
	if err := kraken.Connect(); err != nil {
		log.Fatalf("Error connecting to web socket: %s", err.Error())
	}

	// subscribe to BTCUSD`s book
	if err := kraken.SubscribeBook([]string{ws.BTCUSD}, ws.Depth10); err != nil {
		log.Fatalf("SubscribeBook error: %s", err.Error())
	}

	// subscribe to BTCUSD`s candles
	if err := kraken.SubscribeCandles([]string{ws.BTCUSD}, ws.Interval1440); err != nil {
		log.Fatalf("SubscribeCandles error: %s", err.Error())
	}

	// subscribe to BTCUSD`s ticker
	if err := kraken.SubscribeTicker([]string{ws.BTCUSD}); err != nil {
		log.Fatalf("SubscribeTicker error: %s", err.Error())
	}

	// subscribe to BTCUSD`s trades
	if err := kraken.SubscribeTrades([]string{ws.BTCUSD}); err != nil {
		log.Fatalf("SubscribeTrades error: %s", err.Error())
	}

	// subscribe to BTCUSD`s spread
	if err := kraken.SubscribeSpread([]string{ws.BTCUSD}); err != nil {
		log.Fatalf("SubscribeSpread error: %s", err.Error())
	}

	orderBook := ws.NewOrderBook(10, 5, 8)

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
			case ws.TickerUpdate:
				log.Printf("----Ticker of %s----", update.Pair)
				log.Printf("Ask: %s with %s", data.Ask.Price.String(), data.Ask.Volume.String())
				log.Printf("Bid: %s with %s", data.Bid.Price.String(), data.Bid.Volume.String())
				log.Printf("Open today: %s | Open last 24 hours: %s", data.Open.Today.String(), data.Open.Last24.String())
			case ws.Candle:
				log.Printf("----Candle of %s----", update.Pair)
				log.Printf("Open: %s", data.Open.String())
				log.Printf("High: %s", data.High.String())
				log.Printf("Low: %s", data.Low.String())
				log.Printf("Close: %s", data.Close.String())
			case []ws.Trade:
				log.Printf("----Trades of %s----", update.Pair)
				for i := range data {
					log.Printf("Price: %s", data[i].Price.String())
					log.Printf("Volume: %s", data[i].Volume.String())
					log.Printf("Time: %s", data[i].Time.String())
					log.Printf("Order type: %s", data[i].OrderType)
					log.Printf("Side: %s", data[i].Side)
					log.Printf("Misc: %s", data[i].Misc)
				}
			case ws.Spread:
				log.Printf("----Spread of %s----", update.Pair)
				log.Printf("Ask: %s with %s", data.Ask.String(), data.AskVolume.String())
				log.Printf("Bid: %s with %s", data.Bid.String(), data.BidVolume.String())
			case ws.OrderBookUpdate:
				if err := orderBook.ApplyUpdate(data, true); err != nil {
					log.Fatal(err)
				}

				log.Print(orderBook.String())
			default:
			}
		}
	}
}
